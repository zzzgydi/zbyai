package rag

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/service/llm"
	"github.com/zzzgydi/zbyai/service/scrape"
	"github.com/zzzgydi/zbyai/service/search"
)

var _ RAGExec = (*RAGLimit)(nil)

// RAGLimit 限制性rag方案
// 1. 先搜索，爬取文章
// 2. 直接用最先爬到的文章进行构建context
// 3. 根据snippet进行排序 TODO
type RAGLimit struct {
	wg        sync.WaitGroup
	mutex     sync.Mutex
	maxToken  int
	ctx       context.Context
	ctxCancel context.CancelFunc
	seqResult model.SearchList
	retResult model.SearchList
	content   string
	setting   *model.ThreadRunSetting
	logger    *slog.Logger
}

func NewRAGLimit(setting *model.ThreadRunSetting, logger *slog.Logger) (*RAGLimit, error) {
	if setting == nil {
		return nil, fmt.Errorf("setting is nil")
	}
	if setting.QueryList == nil || len(setting.QueryList) == 0 {
		return nil, fmt.Errorf("query list is nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger = logger.With("rag", "limit")

	return &RAGLimit{
		ctx:       ctx,
		ctxCancel: cancel,
		maxToken:  8000, // TODO：根据model不同
		setting:   setting,
		logger:    logger,
	}, nil
}

// Run implements RAGable.
func (r *RAGLimit) Run() (model.SearchList, error) {
	r.logger.Info("Start search run", "setting", r.setting)

	var searchMu sync.Mutex
	var searchWg sync.WaitGroup

	searchResults := make([]*model.ThreadSearch, 0)
	cache := map[string]bool{}

	for idx, query := range r.setting.QueryList {
		searchWg.Add(1)
		go func(idx int, query string) {
			defer searchWg.Done()

			var engine search.Searcher
			if idx == 0 {
				engine = search.NewSearXNG()
			} else {
				engine = search.NewDDG()
			}

			result, err := engine.Search(&search.SearchOption{Query: query, Count: 6})
			if err != nil {
				r.logger.Error("search error", "engine", engine.Engine(), "error", err)
				return
			}

			items := make([]*model.ThreadSearch, len(result.Items))
			for idx, item := range result.Items {
				items[idx] = &model.ThreadSearch{
					Id:        xid.New().String(),
					Title:     item.Title,
					Link:      item.Link,
					Snippet:   item.Snippet,
					Page:      "",
					Token:     0,
					CreatedAt: time.Now(),
				}
			}
			searchMu.Lock()
			searchResults = append(searchResults, items...)
			searchMu.Unlock()

			for _, item := range items {
				r.wg.Add(1)
				go func(item *model.ThreadSearch) {
					defer r.wg.Done()

					// 避免重复爬取
					searchMu.Lock()
					if _, ok := cache[item.Link]; ok {
						searchMu.Unlock()
						return
					}
					cache[item.Link] = true
					searchMu.Unlock()

					page, err := scrape.ScrapeByURL(item.Link)
					if err != nil {
						r.logger.Error("scrape error", "engine", engine.Engine(), "link", item.Link, "error", err)
						return
					}

					item.Page = page
					item.Token = llm.CountTokenText("gpt-3.5-turbo", page)

					r.logger.Info("scrape success", "engine", engine.Engine(), "link", item.Link, "token", item.Token)

					r.checkDone(item)
				}(item)
			}

		}(idx, query)
	}

	searchWg.Wait()

	// 去除重复
	unique := make(map[string]bool)
	retItems := make(model.SearchList, 0)
	for _, item := range searchResults {
		if _, ok := unique[item.Link]; ok {
			continue
		}
		unique[item.Link] = true
		retItems = append(retItems, item)
	}
	r.retResult = retItems

	go func() {
		r.wg.Wait()
		r.ctxCancel()
		r.logger.Info("Search done", "count", len(retItems))
	}()

	return retItems, nil
}

// WaitContent implements RAGable.
func (r *RAGLimit) WaitContent() string {
	select {
	case <-r.ctx.Done():
		break
	case <-time.After(10 * time.Second):
		r.logger.Warn("SearchRun use content by timeout")
		break
	}

	if r.content != "" {
		return r.content
	}

	content := ""
	for idx, item := range r.retResult {
		page := item.Page
		if page == "" {
			page = item.Snippet
		}
		temp := content + fmt.Sprintf("## Page %d\n%s\n\n", idx+1, page)
		if (llm.CountTokenText("gpt-3.5-turbo", temp)) >= r.maxToken {
			continue
		}
		content = temp
	}
	return content
}

// WaitResult implements RAGable.
func (r *RAGLimit) WaitResult() model.SearchList {
	r.wg.Wait()
	return r.retResult
}

// WaitResultUnlock implements RAGable.
func (r *RAGLimit) WaitResultUnlock() {
}

func (r *RAGLimit) checkDone(item *model.ThreadSearch) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.seqResult = append(r.seqResult, item)

	if r.content != "" {
		return
	}

	maybeDone := false
	content := ""

	for idx, item := range r.seqResult {
		if item.Page == "" {
			continue
		}

		temp := content + fmt.Sprintf("## Page %d\n%s\n\n", idx+1, item.Page)

		if (llm.CountTokenText("gpt-3.5-turbo", temp)) >= r.maxToken {
			maybeDone = true
			break // TODO: 其实这里是有问题的
		}
		content = temp
	}

	if maybeDone {
		r.content = content
		r.ctxCancel()
	}
}
