package thread

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common"
	L "github.com/zzzgydi/zbyai/common/logger"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/model"
	"gorm.io/gorm"
)

// 创建thread
// - 查找是否有相同的query，如果有就直接返回id
// - 如果没有就创建一个新的thread
// - 创建一个id，返回id
// - 在redis中保存任务
// - go 执行任务
// -    不断将内容写入redis，并设置过期时间
func CreateThread(userId, query string, logger *slog.Logger) (*model.Thread, *model.ThreadRun, error) {
	var thread *model.Thread
	var threadRun *model.ThreadRun

	query = strings.TrimSpace(query)

	err := common.MDB.Transaction(func(tx *gorm.DB) error {
		title := utils.Ellipsis(query, 100)
		thread = model.NewThread(userId, title, true)
		if err := thread.Create(tx); err != nil {
			return err
		}
		threadRun = model.NewThreadRun(thread.Id, query)
		if err := threadRun.Create(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	if thread == nil || threadRun == nil {
		return nil, nil, errors.New("create thread unhandled error")
	}

	runner, err := NewThreadRunner(thread, nil, logger)
	if err != nil {
		return nil, nil, err
	}

	if err := runner.Run(threadRun); err != nil {
		return nil, nil, err
	}

	return thread, threadRun, nil
}

// append
// 判断当前thread是不是自己的
// 如果不是自己的就fork一个，然后再执行
func AppendThread(userId, id, query string, logger *slog.Logger) (*model.Thread, *model.ThreadRun, error) {
	var thread *model.Thread
	var threadRun *model.ThreadRun
	var runHistory []*model.ThreadRun

	query = strings.TrimSpace(query)

	err := common.MDB.Transaction(func(tx *gorm.DB) error {
		// find the thread
		thread = &model.Thread{}
		if err := tx.Where("id = ?", id).First(thread).Error; err != nil {
			return err
		}

		// check if thread is visible or belongs to user
		if thread.UserId != userId && !thread.Visible {
			return errors.New("thread not found")
		}

		// find all runs
		runs, err := thread.AllRuns(tx)
		if err != nil {
			return err
		}

		// check if belongs to user
		// if not, fork it
		if thread.UserId != userId {
			// fork a new thread
			thread = model.ForkThread(thread, userId, true)
			if err := thread.Create(tx); err != nil {
				return err
			}
			// new run
			threadRun = model.NewThreadRun(thread.Id, query)
			// copy runs
			runHistory = make([]*model.ThreadRun, len(runs))
			for idx, run := range runs {
				runHistory[idx] = run.CopyTo(thread.Id)
			}
			// create all fork runs
			// save the history and the new run
			temp := append(runHistory, threadRun)
			if err := tx.Create(&temp).Error; err != nil {
				return err
			}
			logger.Info("fork thread",
				"user-change", userId+"<->"+thread.UserId,
				"id-change", id+"<->"+thread.Id)
		} else {
			// new run
			threadRun = model.NewThreadRun(thread.Id, query)
			// set history
			runHistory = runs
			// only create the new run
			if err := threadRun.Create(tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	runner, _ := NewThreadRunner(thread, runHistory, logger)

	if err := runner.Run(threadRun); err != nil {
		return nil, nil, err
	}
	return thread, threadRun, nil
}

// rerun
// 重新执行一些run
// 可以是当前的run，也可以是之前的run
// 执行llm回答，或者别的
// 同样是main的key的话，只能执行一个
func RewriteThread(userId, id string, runId uint64, logger *slog.Logger) (*model.Thread, *model.ThreadRun, error) {
	var thread *model.Thread
	var threadRun *model.ThreadRun
	var threadRunList []*model.ThreadRun

	err := common.MDB.Transaction(func(tx *gorm.DB) error {
		// find the thread
		// rerun 的话，必须是自己的thread
		thread = &model.Thread{}
		if err := tx.Where("id = ? AND user_id = ?", id, userId).First(thread).Error; err != nil {
			return err
		}

		// find all runs
		runs, err := thread.AllRuns(tx)
		if err != nil {
			return err
		}

		threadRunList = runs

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// 历史到当前run，后面的不要
	threadRunBefore := make([]*model.ThreadRun, 0)
	// find the run
	for _, run := range threadRunList {
		if run.Id == runId {
			threadRun = run
			break
		} else {
			threadRunBefore = append(threadRunBefore, run)
		}
	}
	if threadRun == nil {
		return nil, nil, errors.New("run not found")
	}

	runner, err := NewThreadRunner(thread, threadRunBefore, logger)
	if err != nil {
		return nil, nil, err
	}

	if err := runner.Rewrite(threadRun); err != nil {
		return nil, nil, err
	}
	return thread, threadRun, nil
}

// 获取thread
// - 传参数id、offset
// - 从redis获取thread，判断是否存在
// - 如果不存在，就去db中查询，如果有，就全部返回（有个reset语义）
// - 如果有offset，就先从redis获取offset之后的内容
// - 如果没有offset，就从头开始传
// - for 需要select一下ctx.Done()
// - 返回内容
func DetailThread(userId, id string) (*ThreadDetail, error) {
	// 根据threadid获取所有的thread run
	db := common.MDB

	thread := &model.Thread{}
	if err := db.Where("id = ?", id).First(thread).Error; err != nil {
		return nil, err
	}
	if thread.UserId != userId && !thread.Visible {
		return nil, errors.New("thread not found")
	}

	threadRuns, err := thread.AllRuns(db)
	if err != nil {
		return nil, err
	}

	for _, run := range threadRuns {
		if err := run.PrefetchSearch(db); err != nil {
			return nil, err
		}
		if err := run.PrefetchAnswer(db); err != nil {
			return nil, err
		}
		// 判断是否超时
		if run.Status == model.RunStatusInit || run.Status == model.RunStatusWork {
			if run.CreatedAt.Add(3 * time.Minute).Before(time.Now()) {
				run.Status = model.RunStatusError
			}
		}
	}

	current := uint64(0)
	status := model.RunStatusDone
	if len(threadRuns) > 0 {
		last := threadRuns[len(threadRuns)-1]
		current = last.Id
		status = last.Status
	}

	return &ThreadDetail{
		Id:      thread.Id,
		Title:   thread.Title,
		History: threadRuns,
		Current: current,
		Status:  status,
	}, nil
}

func StreamThread(id string, runId uint64, offset int, ctx *gin.Context, logger *slog.Logger) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	if logger == nil {
		logger = slog.New(L.Handler)
	}

	key := keyForStream(id, runId)

	exists, err := common.RDB.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("redis exists", "error", err)
		return
	}
	if exists == 0 {
		return // key not found, maybe done
	}

	for {
		select {
		case <-ctx.Request.Context().Done():
			return
		default:
			// pass
		}

		tasks, err := common.RDB.LRange(ctx, key, int64(offset), -1).Result()
		if err != nil {
			logger.Error("failed to range redis", "error", err)
			break
		}

		for _, task := range tasks {
			ctx.Writer.Write([]byte("data: " + task + "\n\n"))
			ctx.Writer.Flush()

			// TODO optimize
			if strings.Contains(task, `"type":"error"`) ||
				strings.Contains(task, `"type":"done"`) ||
				strings.Contains(task, `"type": "error"`) ||
				strings.Contains(task, `"type": "done"`) {
				return
			}
		}

		offset += len(tasks)

		if offset == 0 {
			time.Sleep(500 * time.Millisecond)
		} else {
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func DeleteThread(userId, id string) error {
	db := common.MDB

	thread := &model.Thread{}
	if err := db.Where("id = ? AND user_id = ?", id, userId).First(thread).Error; err != nil {
		return err
	}

	if err := db.Delete(thread).Error; err != nil {
		return err
	}
	return nil
}

func ListThread(userId string) ([]*model.Thread, error) {
	db := common.MDB

	threads := make([]*model.Thread, 0)
	if err := db.Where("user_id = ?", userId).
		Order("created_at desc").
		Find(&threads).Error; err != nil {
		return nil, err
	}

	return threads, nil
}
