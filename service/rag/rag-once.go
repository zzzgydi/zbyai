package rag

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/service/llm"
	"github.com/zzzgydi/zbyai/service/scrape"
	"github.com/zzzgydi/zbyai/service/search"
)

var _ RAGExec = (*RAGOnce)(nil)

// RAGOnce 只搜一次
type RAGOnce struct {
	maxToken     int
	ctx          context.Context
	ctxCancel    context.CancelFunc
	searchResult model.SearchList
	scrapeWg     sync.WaitGroup
	setting      *model.ThreadRunSetting
	logger       *slog.Logger
	onSearch     func([]*model.ThreadSearch)
}

type innerChunk struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

func NewRAGOnce(setting *model.ThreadRunSetting, onSearch func([]*model.ThreadSearch), logger *slog.Logger) (*RAGOnce, error) {
	if setting == nil {
		return nil, fmt.Errorf("setting is nil")
	}
	if setting.QueryList == nil || len(setting.QueryList) == 0 {
		return nil, fmt.Errorf("query list is nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &RAGOnce{
		ctx:       ctx,
		ctxCancel: cancel,
		maxToken:  8000, // TODO：根据model不同
		setting:   setting,
		logger:    logger,
		onSearch:  onSearch,
	}, nil
}

func (r *RAGOnce) Run() (model.SearchList, error) {
	r.runSearch()
	go func() {
		r.scrapeWg.Wait()
		r.ctxCancel()
	}()
	return r.searchResult, nil
}

func (r *RAGOnce) WaitContent() string {
	select {
	case <-r.ctx.Done():
		break
	case <-time.After(5 * time.Second):
		r.logger.Warn("[RAG] use content by timeout")
		r.ctxCancel()
		break
	}

	return SimpleRerunk(r.searchResult, r.maxToken)
}

func (r *RAGOnce) WaitResult() model.SearchList {
	r.scrapeWg.Wait()
	return r.searchResult
}

func (r *RAGOnce) WaitResultUnlock() {
}

func (r *RAGOnce) runSearch() {
	maxSize := min(3, len(r.setting.QueryList))
	engineList := onceEngineList(maxSize)
	results := make([]*model.ThreadSearch, 0)
	urlMap := make(map[string]bool)

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < maxSize; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			engine := engineList[idx]
			query := r.setting.QueryList[idx]

			count := 4
			if maxSize == 1 {
				count = 6
			}

			result, err := engine.Search(&search.SearchOption{Query: query, Count: count})
			if err != nil {
				r.logger.Error("[RAG] search error", "engine", engine.Engine(), "error", err)
				return
			}

			logResults := []map[string]string{}
			for _, item := range result.Items {
				logResults = append(logResults, map[string]string{
					"title": item.Title,
					"link":  item.Link,
				})
			}
			r.logger.Info("[RAG] search success", "engine", engine.Engine(), "query", query, "result", logResults)

			mutex.Lock() // 保证线程安全

			items := make([]*model.ThreadSearch, 0, len(result.Items))
			for _, item := range result.Items {
				if _, ok := urlMap[item.Link]; ok {
					continue
				}
				urlMap[item.Link] = true
				items = append(items, &model.ThreadSearch{
					Id:        xid.New().String(),
					Title:     item.Title,
					Link:      item.Link,
					Snippet:   item.Snippet,
					Page:      "",
					Token:     0,
					CreatedAt: time.Now(),
				})
			}

			results = append(results, items...)
			mutex.Unlock()

			if r.onSearch != nil {
				r.onSearch(items)
			}
		}(i)
	}

	wg.Wait()

	r.searchResult = results
	for _, item := range results {
		go r.runScrape(item)
	}
}

func (r *RAGOnce) runScrape(item *model.ThreadSearch) {
	if r.ctx.Err() != nil {
		return
	}
	r.scrapeWg.Add(1)
	defer r.scrapeWg.Done()

	scrapeClient := scrape.NewScrape(true)
	scrapeClient.AddPipeline(utils.RemoveMarkdownImages)

	content, err := scrapeClient.Run(r.ctx, item.Link)
	if err != nil {
		r.logger.Error("[RAG] scrape error", "link", item.Link, "error", err)
		return
	}

	item.Page = content
	item.Token = llm.CountTokenText("gpt-3.5-turbo", content)
	r.logger.Info("[RAG] scrape success", "link", item.Link, "token", item.Token)
}
