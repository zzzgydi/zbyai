package thread

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/model"
)

type Streamable interface {
	Stream() string
}

type ThreadStream struct {
	count    uint32
	logger   *slog.Logger
	threadId string
	runId    uint64
}

func NewThreadStream(threadId string, runId uint64, logger *slog.Logger) *ThreadStream {
	ctx := context.Background()
	key := keyForStream(threadId, runId)
	// 判断redis中是否有这个key，如果存在就清空
	if common.RDB.Exists(ctx, key).Val() == 1 {
		common.RDB.Del(ctx, key)
	}

	return &ThreadStream{
		count:    0,
		logger:   logger,
		threadId: threadId,
		runId:    runId,
	}
}

func (ts *ThreadStream) SendStart(query string) {
	ts.send("query", "")
}

func (ts *ThreadStream) SendSetting(setting *model.ThreadRunSetting) {
	ts.send("setting", setting)
}

func (ts *ThreadStream) SendSearchBegin() {
	ts.send("search", map[string]any{
		"status": 1, // 开始搜索
	})
}

func (ts *ThreadStream) SendSearchDone(search model.SearchList) {
	ts.send("search", map[string]any{
		"status": 2, // 搜索完成，开始爬取
		"search": search,
	})
}

func (ts *ThreadStream) SendSearchError(err error) {
	ts.send("search", map[string]any{
		"status": 3,
		"errMsg": err.Error(),
	})
}

func (ts *ThreadStream) SendScrapeDone() {
	ts.send("scrape", map[string]any{
		"status": 2,
	})
}

func (ts *ThreadStream) SendAnswerBegin(answer *model.ThreadAnswer) {
	answer.Status = model.AnswerInit
	ts.send("answer", map[string]any{
		"id":     answer.Id,
		"key":    answer.Key,
		"status": answer.Status,
		"model":  answer.Model,
	})
}

func (ts *ThreadStream) SendAnswerDelta(answer *model.ThreadAnswer, delta string) {
	answer.Status = model.AnswerWork
	ts.send("answer", map[string]any{
		"id":     answer.Id,
		"key":    answer.Key,
		"status": answer.Status,
		"delta":  delta,
	})
}

func (ts *ThreadStream) SendAnswerDone(answer *model.ThreadAnswer) {
	answer.Status = model.AnswerDone
	ts.send("answer", map[string]any{
		"id":     answer.Id,
		"key":    answer.Key,
		"status": answer.Status,
	})
}

func (ts *ThreadStream) SendAnswerError(answer *model.ThreadAnswer, err error) {
	answer.Status = model.AnswerError
	ts.send("answer", map[string]any{
		"id":     answer.Id,
		"key":    answer.Key,
		"status": answer.Status,
		"errMsg": err.Error(),
	})
}

func (ts *ThreadStream) SendError(err error) {
	ts.send("error", map[string]any{
		"errMsg": err.Error(),
	})
}

func (ts *ThreadStream) SendDone() {
	ts.send("done", nil)
}

// Release will delete the stream from redis
func (ts *ThreadStream) Release() error {
	key := keyForStream(ts.threadId, ts.runId)
	// 设置10s后过期
	return common.RDB.Expire(context.Background(), key, 10*time.Second).Err()
}

func (ts *ThreadStream) send(typ string, data any) {
	atomic.AddUint32(&ts.count, 1)

	sendData := map[string]any{
		"id":   ts.count,
		"type": typ,
		"data": data,
	}
	data, err := json.Marshal(sendData)
	if err != nil {
		ts.logger.Error("json marshal error", "error", err)
		return
	}

	key := keyForStream(ts.threadId, ts.runId)
	err = common.RDB.RPush(context.Background(), key, data).Err()
	if err != nil {
		ts.logger.Error("redis rpush error", "error", err)
		return
	}

	// 5mins
	err = common.RDB.Expire(context.Background(), key, 5*time.Minute).Err()
	if err != nil {
		ts.logger.Error("redis expire error", "error", err)
	}
}
