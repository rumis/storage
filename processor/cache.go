package processor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/locker"
	"github.com/rumis/storage/v2/meta"
	"github.com/rumis/storage/v2/scache"
)

// RedisReadProcessor 缓存读取处理
type RedisReadProcessor struct {
	locker locker.Locker
	reader meta.KeyValueReader
	writer meta.KeyValueWriter
	next   ReadProcessor
}

// Read 数据读取
func (r *RedisReadProcessor) Read(ctx context.Context, in interface{}, out interface{}, exp time.Duration) error {
	zero, ok := out.(meta.Zero)
	if !ok {
		return errors.New("params out must implements Zero interface")
	}
	// 读取缓存
	err := r.reader(ctx, in, out)
	if err == nil {
		// 缓存读取成功，直接返回
		return nil
	}
	if err != nil && err != redis.Nil {
		// 缓存不存在或读取错误
		fmt.Println("记录错误")
	}
	// 读取错误&缓存中key不存在都继续执行以下流程
	// 锁
	keyIT, ok := in.(meta.Key)
	if !ok {
		return meta.EI_KeyNotImplement
	}
	// 不再包含后续操作，直接返回
	if r.next == nil {
		return nil
	}

	// 后续操作先尝试下加锁
	ok = r.locker.Adder(ctx, keyIT.Key())
	if !ok {
		// 未抢到锁 - 尝试多次读取缓存
		for i := 0; i < r.locker.RetryTimes; i++ {
			time.Sleep(r.locker.RetrySpan)
			err = r.reader(ctx, in, out)
			if err == nil {
				// 数据从缓存中读取成功，直接返回
				return nil
			}
		}
	}
	// 缓存未读到数据 继续执行下一流程
	err = r.next.Read(ctx, in, out, exp)
	if err != nil {
		// 失败，返回错误
		return err
	}
	// 如果库中不存在数据
	if zero.Zero() {
		// 写入个空数据
		exp = r.locker.Expire
	}
	// 写缓存
	err = r.Write(ctx, out, exp)
	if err != nil {
		fmt.Println("redis write erro")
	}
	// 删除锁
	r.locker.Deleter(ctx, keyIT.Key())

	return nil
}

// Write 数据写入
func (r *RedisReadProcessor) Write(ctx context.Context, in interface{}, exp time.Duration) error {
	return r.writer(ctx, in, exp)
}

// NewRedisCacheReaderProcessor 新建Redis缓存处理器
func NewRedisCacheReaderProcessor(next ReadProcessor) ReadProcessor {
	return &RedisReadProcessor{
		next:   next,
		locker: locker.DefaultLocker(),
		reader: scache.NewRedisKeyValueReader(scache.WithClient(scache.DefaultClient())),
		writer: scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient())),
	}
}
