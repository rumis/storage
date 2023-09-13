package processor

import (
	"context"

	"github.com/rumis/storage/v2/locker"
	"github.com/rumis/storage/v2/meta"
)

// RedisReadProcessor 缓存读取处理
type RedisReadProcessor struct {
	locker locker.Locker
	reader meta.KeyValueReader
	writer meta.KeyValueWriter
	next   ReadProcessor
}

// Read 数据读取
func (r *RedisReadProcessor) Read(ctx context.Context, in interface{}, out interface{}) bool {

	return false
}

// Write 数据写入
func (r *RedisReadProcessor) Write(ctx context.Context, in interface{}, out interface{}) bool {

	return false
}

// Lock 加锁
func (r *RedisReadProcessor) Lock(ctx context.Context) error {
	if r.locker.Adder == nil || r.locker.Deleter == nil {
		return nil
	}
	r.locker.Adder(ctx, "")
	return nil
}

// UnLock 释放锁
func (r *RedisReadProcessor) UnLock(ctx context.Context) error {
	if r.locker.Deleter == nil {
		return nil
	}
	return r.locker.Deleter(ctx, "")
}

// NewRedisCacheReaderProcessor 新建Redis缓存处理器
func NewRedisCacheReaderProcessor(next ReadProcessor) ReadProcessor {

	return &RedisReadProcessor{}
}
