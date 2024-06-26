package rag

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/rs/xid"
	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/service/llm"
	"github.com/zzzgydi/zbyai/service/scrape"
	"github.com/zzzgydi/zbyai/service/search"
)

var _ RAGExec = (*RAGCommon)(nil)

// 常规性rag方案
// 1. 搜索，然后爬取
// 2. 分块，embedding之后排序
// 3. 取topN个块，拼接成content
type RAGCommon struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	maxToken  int
	setting   *model.ThreadRunSetting
	logger    *slog.Logger

	// query
	queryList []string
	queryEmbs []openai.Embedding

	// search
	searchWg     sync.WaitGroup
	searchMutex  sync.Mutex
	searchResult model.SearchList

	// scrape
	scrapeWg sync.WaitGroup

	// cache
	cacheMap   map[string]bool
	cacheMutex sync.Mutex

	// chunk
	chunkList  []*RAGChunk
	chunkMutex sync.Mutex
}

func NewRAGCommon(setting *model.ThreadRunSetting, logger *slog.Logger) *RAGCommon {
	if setting == nil {
		logger.Warn("Setting is nil")
		return nil
	}
	if setting.QueryList == nil || len(setting.QueryList) == 0 {
		logger.Warn("QueryList is nil")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	queryList := setting.QueryList

	if len(setting.QueryList) > 3 {
		queryList = setting.QueryList[:3]
	}

	logger = logger.With("rag", "common")

	return &RAGCommon{
		ctx:       ctx,
		ctxCancel: cancel,
		maxToken:  6000, // TODO：根据model不同
		setting:   setting,
		queryList: queryList,
		logger:    logger,
		cacheMap:  make(map[string]bool),
	}
}

func (r *RAGCommon) Run() (model.SearchList, error) {
	go r.runQueryEmb()
	r.runSearch()

	r.searchMutex.Lock()
	defer r.searchMutex.Unlock()

	// 去除重复
	unique := make(map[string]bool)
	retItems := make(model.SearchList, 0)
	for _, item := range r.searchResult {
		if _, ok := unique[item.Link]; ok {
			continue
		}
		unique[item.Link] = true
		retItems = append(retItems, item)
	}

	return retItems, nil
}

func (r *RAGCommon) WaitContent() string {
	select {
	case <-r.ctx.Done():
		break
	case <-time.After(10 * time.Second):
		r.logger.Warn("[RAG] use content by timeout")
		r.ctxCancel()
		break
	}

	// 计算max chunk
	r.chunkMutex.Lock()
	defer r.chunkMutex.Unlock()

	for _, queryEmb := range r.queryEmbs {
		for _, chunk := range r.chunkList {
			score, err := queryEmb.DotProduct(&chunk.Embedding)
			if err != nil {
				r.logger.Error("[RAG] dot product error", "error", err)
				continue
			}
			chunk.Score += score
		}
	}

	// 降序，所以这里用大于号
	sort.Slice(r.chunkList, func(i, j int) bool {
		return r.chunkList[i].Score > r.chunkList[j].Score
	})

	for idx, query := range r.queryList {
		r.logger.Info("[RAG] query", "idx", idx, "text", query)
	}
	for idx, item := range r.chunkList {
		r.logger.Info("[RAG] chunk", "idx", idx, "score", item.Score, "text", item.Text)
	}

	// 生成结果
	content := ""

	for idx, item := range r.chunkList {
		temp := content + fmt.Sprintf("## Doc %d\n%s\n\n", idx+1, item.Text)
		if (llm.CountTokenText("gpt-3.5-turbo", temp)) >= r.maxToken {
			break
		}
		content = temp
	}

	return content
}

func (r *RAGCommon) WaitResult() model.SearchList {
	// 等待所有爬取完成后，保存结果
	r.scrapeWg.Wait()
	r.searchMutex.Lock()

	// 去除重复
	unique := make(map[string]bool)
	retItems := make(model.SearchList, 0)
	for _, item := range r.searchResult {
		if _, ok := unique[item.Link]; ok {
			continue
		}
		unique[item.Link] = true
		retItems = append(retItems, item)
	}

	return retItems
}

func (r *RAGCommon) WaitResultUnlock() {
	r.searchMutex.Unlock()
}

func (r *RAGCommon) runQueryEmb() {
	r.logger.Info("[RAG] start query embedding run", "setting", r.setting)
	r.searchWg.Add(1)
	defer r.searchWg.Done()

	for idx := 0; idx < 3; idx++ {
		embList, err := llm.QueryEmbedding(r.ctx, r.queryList)
		if err != nil {
			r.logger.Error("[RAG] query embedding error", "retry", idx+1, "error", err)
			continue
		}
		r.queryEmbs = embList
		break
	}
}

func (r *RAGCommon) runSearch() {
	r.logger.Info("[RAG] start search run", "setting", r.setting)

	for idx, query := range r.queryList {
		r.searchWg.Add(1)

		go func(idx int, query string) {
			defer r.searchWg.Done()

			if r.ctx.Err() != nil {
				return
			}

			var engine search.Searcher
			if idx == 0 {
				engine = search.NewSearXNG()
			} else {
				engine = search.NewDDG()
			}

			result, err := engine.Search(&search.SearchOption{Query: query, Count: 6})
			if err != nil {
				r.logger.Error("[RAG] search error", "engine", engine.Engine(), "error", err)
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

			r.searchMutex.Lock()
			r.searchResult = append(r.searchResult, items...)
			r.searchMutex.Unlock()

			for _, item := range items {
				go r.runScrape(engine.Engine(), item)
			}
		}(idx, query)
	}

	// 等所有搜索的完成
	r.searchWg.Wait()
}

func (r *RAGCommon) runScrape(engine string, item *model.ThreadSearch) {
	if r.ctx.Err() != nil {
		return
	}

	// 避免重复爬取
	r.cacheMutex.Lock()
	if _, ok := r.cacheMap[item.Link]; ok {
		r.cacheMutex.Unlock()
		return
	}
	r.cacheMap[item.Link] = true
	r.cacheMutex.Unlock()

	r.scrapeWg.Add(1)
	defer r.scrapeWg.Done()

	page, err := scrape.ScrapeByURL(item.Link)
	if err != nil {
		r.logger.Error("[RAG] scrape error", "engine", engine, "link", item.Link, "error", err)
		return
	}

	isCode := r.setting != nil && r.setting.IsProgramming
	if isCode {
		// 去掉所有链接部分
		page = utils.RemoveMarkdownLink(page)
	}

	item.Page = page
	item.Token = llm.CountTokenText("gpt-3.5-turbo", page)

	r.logger.Info("[RAG] scrape success", "engine", engine, "link", item.Link, "token", item.Token)

	if err := r.runChunks(page); err != nil {
		r.logger.Error("[RAG] run chunks error", "error", err)
	}
}

func (r *RAGCommon) runChunks(page string) error {
	if r.ctx.Err() != nil {
		return nil
	}

	mdSplitter := textsplitter.NewMarkdownTextSplitter(
		textsplitter.WithChunkSize(1280),
		textsplitter.WithChunkOverlap(128),
	)

	docs, err := mdSplitter.SplitText(page)
	if err != nil {
		return err
	}

	// 过滤一些不合适的文本
	temp := make([]string, 0)
	for _, doc := range docs {
		doc = strings.TrimSpace(doc)
		if doc == "" || utf8.RuneCountInString(doc) < 33 {
			continue
		}
		temp = append(temp, doc)
	}

	// 仅保留中间的20个
	maxDocs := 20
	if len(temp) > maxDocs {
		l := len(temp) - maxDocs
		temp = temp[l/2 : l/2+maxDocs]
	}

	var wg sync.WaitGroup

	docGroups := groupText(temp, 5)

	for _, group := range docGroups {
		wg.Add(1)
		go func(group []string) {
			defer wg.Done()

			if r.ctx.Err() != nil {
				return
			}
			embList, err := llm.AnswerEmbedding(r.ctx, group)
			if err != nil {
				r.logger.Error("[RAG] text embedding error", "error", err)
				return
			}

			r.chunkMutex.Lock()
			defer r.chunkMutex.Unlock()

			for i, emb := range embList {
				r.chunkList = append(r.chunkList, &RAGChunk{
					Text:      group[i],
					Embedding: emb,
				})
			}
		}(group)
	}

	wg.Wait()

	return nil
}
