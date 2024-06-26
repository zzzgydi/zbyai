package thread

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"time"

	"github.com/rs/xid"
	"github.com/sashabaranov/go-openai"
	"github.com/zzzgydi/zbyai/common"
	L "github.com/zzzgydi/zbyai/common/logger"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/model"
	R "github.com/zzzgydi/zbyai/service/rag"
	"gorm.io/gorm"
)

type ThreadRunner struct {
	db          *gorm.DB
	thread      *model.Thread
	current     *model.ThreadRun
	history     []*model.ThreadRun
	chatHistory []openai.ChatCompletionMessage
	logger      *slog.Logger
	rag         R.RAGExec
	ragContent  string
	stream      *ThreadStream
}

func NewThreadRunner(thread *model.Thread, history []*model.ThreadRun, logger *slog.Logger) (*ThreadRunner, error) {
	if logger == nil {
		logger = slog.New(L.Handler)
	}
	logger = logger.With("thread", thread.Id)

	return &ThreadRunner{
		db:      common.MDB,
		thread:  thread,
		history: history,
		logger:  logger,
	}, nil
}

func (t *ThreadRunner) Run(run *model.ThreadRun) error {
	t.current = run

	// lock to limit append thread
	if err := LockThread(t.thread.Id); err != nil {
		return err
	}

	t.logger = t.logger.With("runId", run.Id)
	t.stream = NewThreadStream(t.thread.Id, run.Id, t.logger)
	t.stream.SendStart("")

	go func() {
		// release lock
		defer UnlockThread(t.thread.Id)

		t.runChatHistory()
		t.runRephrased()
		t.runSearch()
		t.runLLM()
		t.done()
	}()

	return nil
}

func (t *ThreadRunner) Rewrite(run *model.ThreadRun) error {
	t.current = run

	// lock to limit append thread
	if err := LockThread(t.thread.Id); err != nil {
		return err
	}

	t.logger = t.logger.With("runId", run.Id)
	t.stream = NewThreadStream(t.thread.Id, run.Id, t.logger)
	t.stream.SendStart("")

	// release lock
	go func() {
		defer UnlockThread(t.thread.Id)
		t.runChatHistory()
		run.PrefetchSearch(t.db)

		// 增加一点随机性
		t.ragContent = R.SimpleRerunk(run.Search, 7000+rand.Intn(2000))

		t.runLLM()
		t.done()
	}()

	return nil
}

func (t *ThreadRunner) runChatHistory() {
	his := t.history
	if len(his) > 4 {
		his = his[len(his)-4:]
	}

	history := make([]openai.ChatCompletionMessage, 0)
	for _, h := range his {
		h.PrefetchAnswer(t.db)
		mainAns := ""
		for _, a := range h.Answer {
			// 一直找，找到最后一个
			if a.Key == model.AnswerKeyMain && a.Status == model.AnswerDone {
				mainAns = a.Content
			}
		}
		if mainAns == "" {
			continue
		}
		history = append(history,
			openai.ChatCompletionMessage{Role: "user", Content: h.Query},
			openai.ChatCompletionMessage{Role: "assistant", Content: mainAns},
		)
	}
	t.chatHistory = history
}

func (t *ThreadRunner) runRephrased() {
	run := t.current

	need, err := peSearchJudge(run.Query, t.chatHistory, t.logger)
	if err != nil {
		t.logger.Error("[Thread] run rephrased error", "error", err)
		run.Setting = &model.ThreadRunSetting{
			UseSearch:     true,
			IsProgramming: false,
			QueryList:     []string{utils.Ellipsis(run.Query, 600)},
		}
	} else {
		queryList := []string{}
		switch rephrased := need.Rephrased.(type) {
		case []any:
			for _, item := range rephrased {
				if str, ok := item.(string); ok {
					queryList = append(queryList, utils.Ellipsis(str, 600))
				}
			}
		default:
			if rephrased != nil {
				queryList = []string{utils.Ellipsis(fmt.Sprintf("%v", rephrased), 1000)}
			}
		}

		run.Setting = &model.ThreadRunSetting{
			UseSearch:     need.UseSearch,
			Model:         need.Model,
			IsProgramming: need.IsProgramming,
			Language:      need.Language,
			QueryList:     queryList,
		}
	}

	t.stream.SendSetting(run.Setting)
}

func (t *ThreadRunner) runSearch() {
	run := t.current
	run.Status = model.RunStatusWork

	onSearch := func(items []*model.ThreadSearch) {
		t.stream.SendSearchDone(items)
	}

	rag, err := R.NewRAGOnce(run.Setting, onSearch, t.logger)
	// rag, err := R.NewRAGLimit(run.Setting, t.logger)
	if err != nil {
		return
	}

	t.rag = rag
	t.stream.SendSearchBegin()

	items, err := t.rag.Run()
	if err != nil || items == nil {
		t.stream.SendSearchError(errors.New("search error"))
	}

	for _, item := range items {
		run.SearchIds = append(run.SearchIds, item.Id) // save search ids
	}

	go func() {
		results := t.rag.WaitResult()
		defer t.rag.WaitResultUnlock()

		if results == nil {
			return
		}

		if err := t.db.Create(results).Error; err != nil {
			t.logger.Error("[Thread] save search error", "error", err)
		}
	}()

	t.ragContent = t.rag.WaitContent()
}

func (t *ThreadRunner) runLLM() {
	run := t.current

	llmClient, messages := peMainRAG(run.Query, t.ragContent, t.chatHistory, t.logger)

	stream, err := llmClient.StreamCompletion(context.Background(), messages)
	if err != nil {
		t.logger.Error("[Thread] llm request error", "model", llmClient.Model, "error", err)
		t.stream.SendError(err)
		run.Status = model.RunStatusError
		return
	}

	answer := &model.ThreadAnswer{
		Id:        xid.New().String(),
		Key:       model.AnswerKeyMain,
		Status:    model.AnswerInit,
		Model:     llmClient.Model,
		Content:   "",
		CreatedAt: time.Now(),
	}
	t.stream.SendAnswerBegin(answer)

	defer func() {
		answer.Status = model.AnswerDone
		run.Answer = append(run.Answer, answer)
		run.AnswerIds = append(run.AnswerIds, answer.Id)

		if err := t.db.Create(answer).Error; err != nil {
			t.logger.Error("[Thread] save answer error", "error", err)
		}
	}()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			t.stream.SendAnswerDone(answer)
			return
		}
		if err != nil {
			t.logger.Error("[Thread] llm stream recv error", "error", err)
			t.stream.SendAnswerError(answer, err)
			run.Status = model.RunStatusError
			return
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta.Content
			answer.Content += delta
			t.stream.SendAnswerDelta(answer, delta)
		} else {
			t.logger.Error("[Thread] llm stream recv empty", "response", response)
		}
	}
}

func (t *ThreadRunner) done() {
	t.stream.SendDone()

	run := t.current

	if run.Status != model.RunStatusError {
		run.Status = model.RunStatusDone
	}
	run.UpdatedAt = time.Now()
	if err := t.db.Save(run).Error; err != nil {
		t.logger.Error("[Thread] save run error", "error", err)
		// return // 下面的东西等到时间了自己释放吧
	}

	// clear redis cache
	if err := t.stream.Release(); err != nil {
		t.logger.Error("[Thread] release stream error", "error", err)
	}
}
