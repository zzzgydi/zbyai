package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zzzgydi/zbyai/common"
)

// RateLimiter 使用Redis实现的 限流器
type RateLimiter struct {
	flag  string // 区分不同目的
	rate  int    // 允许的最大请求数
	burst int    // 时间窗口大小，单位秒
}

func NewRateLimiter(flag string, rate, burst int) *RateLimiter {
	return &RateLimiter{
		flag:  flag,
		rate:  rate,
		burst: burst,
	}
}

// Allow 检查内容是否被允许请求
func (l *RateLimiter) Allow(ctx context.Context, value string) (bool, error) {
	key := fmt.Sprintf("%s:%s", l.flag, value)
	now := time.Now().UnixNano()

	// 使用pipeline优化Redis命令执行
	pipe := common.RDB.Pipeline()

	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(l.burst)*1e9))
	zcard := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, time.Duration(l.burst)*time.Second)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := zcard.Result()
	if err != nil {
		return false, err
	}
	return count <= int64(l.rate), nil
}
