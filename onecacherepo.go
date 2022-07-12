package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/locker"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/pkg/ujson"
	"github.com/rumis/storage/scache"
	"github.com/rumis/storage/srepo"
)

// NewOneCacheReader 缓存对象读取
func NewOneCacheReader(prefix string) scache.RedisKeyValueObjectReader {
	r := scache.NewRedisKeyValueStringReader(scache.WithClient(scache.DefaultClient()), scache.WithPrefix(prefix))
	return func(ctx context.Context, params interface{}, out interface{}) error {
		res, err := r(ctx, fmt.Sprint(params))
		if err != nil {
			return err
		}
		resp, ok := res.(string)
		if !ok {
			return errors.New("format error")
		}
		err = ujson.Unmarshal([]byte(resp), out)
		if err != nil {
			return err
		}
		return nil
	}
}

// NewOneCacheWriter 缓存对象写入
func NewOneCacheWriter(prefix string) func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
	w := scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()), scache.WithPrefix(prefix))
	return func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
		err := w(ctx, kv, expire)
		return err
	}
}

// NewOneRepoReader 通用数据库读取
func NewOneRepoReader(tableName string, columns []string) func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
	r := srepo.NewSealMysqlOneReader(srepo.WithName(tableName), srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns))
	return func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
		err := r(ctx, out, hand...)
		return err
	}
}

// NewDefaultLocker 默认通用锁
func NewDefaultLocker(biz string) locker.Locker {
	return scache.DefaultRedisLocker(scache.DefaultClient(), biz)
}

// NewOneCacheRepoReader 通用缓存-库数据读取器,单对象
func NewOneCacheRepoReader(prefix string, tablename string, columns []string, biz string) func(ctx context.Context, params interface{}, expire time.Duration, out interface{}, opts ...srepo.ClauseHandler) error {
	cacheReader := NewOneCacheReader(prefix)
	cacheWriter := NewOneCacheWriter(prefix)
	repoReader := NewOneRepoReader(tablename, columns)
	locker := NewDefaultLocker(biz)
	return func(ctx context.Context, params interface{}, expire time.Duration, out interface{}, opts ...srepo.ClauseHandler) error {
		zero, ok := out.(meta.Zero)
		if !ok {
			return errors.New("params out must implements Zero interface")
		}
		// 读取缓存
		err := cacheReader(ctx, params, out)
		if err == nil {
			// 缓存读取成功，直接返回
			return nil
		}
		if err != nil && err != redis.Nil {
			// 缓存读取错误
			fmt.Println("记录错误")
		}
		// 读取错误&缓存中key不存在都继续执行以下流程

		// 锁
		l, err := locker.Reader(ctx, fmt.Sprint(params))
		if err == nil && l != "" {
			// 未抢到锁
			for i := 0; i < locker.RetryTimes; i++ {
				time.Sleep(locker.RetrySpan)
				err = cacheReader(ctx, params, out)
				if err == nil {
					// 数据从缓存中读取成功，直接返回
					return nil
				}
			}
		}
		// 更新锁，读库
		err = locker.Writer(ctx, fmt.Sprint(params))
		if err != nil {
			fmt.Println("记录错误", err)
		}
		// 读库
		err = repoReader(ctx, out, opts...)
		if err != nil {
			// 读库失败，返回错误
			return err
		}
		if zero.Zero() {
			// 写入个空数据
			expire = locker.Expire
		}
		// 写缓存
		buf, err := ujson.Marshal(out)
		if err != nil {
			return err
		}
		err = cacheWriter(ctx, scache.Pair{
			Key:   fmt.Sprint(params),
			Value: string(buf),
		}, expire)
		if err != nil {
			fmt.Println("redis write erro")
		}
		// 删除锁
		locker.Deleter(ctx, fmt.Sprint(params))

		return nil
	}
}
