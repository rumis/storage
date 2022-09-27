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

// OneCacheRepoOptionsHandler 单一对象缓存配置处理方法
type OneCacheRepoOptionsHandler func(*OneCacheRepoOptions)

// OneCacheRepoOptions 单一对象缓存配置
type OneCacheRepoOptions struct {
	CacheReader scache.RedisKeyValueObjectReader
	CacheWriter scache.RedisKeyValueWriter
	RepoReader  srepo.RepoGroupReader
	Locker      locker.Locker
}

// NewOneCacheRepoOptions 创建新的单一对象缓存配置
func NewOneCacheRepoOptions(hand ...OneCacheRepoOptionsHandler) OneCacheRepoOptions {
	opts := OneCacheRepoOptions{}
	for _, h := range hand {
		h(&opts)
	}
	return opts
}

// WithCacheReader 缓存读取器
func WithCacheReader(r scache.RedisKeyValueObjectReader) OneCacheRepoOptionsHandler {
	return func(opts *OneCacheRepoOptions) {
		opts.CacheReader = r
	}
}

// WithCacheWriter 缓存写入
func WithCacheWriter(w scache.RedisKeyValueWriter) OneCacheRepoOptionsHandler {
	return func(opts *OneCacheRepoOptions) {
		opts.CacheWriter = w
	}
}

// WithRepoReader 数据库读取
func WithRepoReader(r srepo.RepoGroupReader) OneCacheRepoOptionsHandler {
	return func(opts *OneCacheRepoOptions) {
		opts.RepoReader = r
	}
}

// WithLocker 锁
func WithLocker(l locker.Locker) OneCacheRepoOptionsHandler {
	return func(opts *OneCacheRepoOptions) {
		opts.Locker = l
	}
}

// // NewOneCacheReader 缓存对象读取
// func NewOneCacheReader(prefix string) scache.RedisKeyValueObjectReader {
// 	r := scache.NewRedisKeyValueStringReader(scache.WithClient(scache.DefaultClient()), scache.WithPrefix(prefix))
// 	return func(ctx context.Context, params interface{}, out interface{}) error {
// 		res, err := r(ctx, fmt.Sprint(params))
// 		if err != nil {
// 			return err
// 		}
// 		resp, ok := res.(string)
// 		if !ok {
// 			return errors.New("format error")
// 		}
// 		err = ujson.Unmarshal([]byte(resp), out)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// // NewOneCacheWriter 缓存对象写入
// func NewOneCacheWriter(prefix string) func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
// 	w := scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()), scache.WithPrefix(prefix))
// 	return func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
// 		err := w(ctx, kv, expire)
// 		return err
// 	}
// }

// // NewOneRepoReader 通用数据库读取
// func NewOneRepoReader(tableName string, columns []string) func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
// 	r := srepo.NewSealMysqlOneReader(srepo.WithName(tableName), srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns))
// 	return func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
// 		err := r(ctx, out, hand...)
// 		return err
// 	}
// }

// // NewDefaultLocker 默认通用锁
// func NewDefaultLocker(biz string) locker.Locker {
// 	return scache.DefaultRedisLocker(scache.DefaultClient(), biz)
// }

// NewOneCacheRepoReader 通用缓存-库数据读取器,单对象
func NewOneCacheRepoReader(opts OneCacheRepoOptions) func(ctx context.Context, params interface{}, expire time.Duration, out interface{}) error {
	return func(ctx context.Context, params interface{}, expire time.Duration, out interface{}) error {
		zero, ok := out.(meta.Zero)
		if !ok {
			return errors.New("params out must implements Zero interface")
		}
		// 读取缓存
		err := opts.CacheReader(ctx, params, out)
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
		ok = opts.Locker.Adder(ctx, fmt.Sprint(params))
		if !ok {
			// 未抢到锁 - 尝试多次读取缓存
			for i := 0; i < opts.Locker.RetryTimes; i++ {
				time.Sleep(opts.Locker.RetrySpan)
				err = opts.CacheReader(ctx, params, out)
				if err == nil {
					// 数据从缓存中读取成功，直接返回
					return nil
				}
			}
		}
		// 缓存未读到数据 读库
		err = opts.RepoReader(ctx, out, params)
		if err != nil {
			// 读库失败，返回错误
			return err
		}
		if zero.Zero() {
			// 写入个空数据
			expire = opts.Locker.Expire
		}
		// 写缓存
		buf, err := ujson.Marshal(out)
		if err != nil {
			return err
		}
		err = opts.CacheWriter(ctx, scache.Pair{
			Key:   fmt.Sprint(params),
			Value: string(buf),
		}, expire)
		if err != nil {
			fmt.Println("redis write erro")
		}
		// 删除锁
		opts.Locker.Deleter(ctx, fmt.Sprint(params))

		return nil
	}
}
