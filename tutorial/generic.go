package tutorial

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/seal"
	"github.com/rumis/seal/query"
	"github.com/rumis/storage/locker"
	"github.com/rumis/storage/pkg/ujson"
	"github.com/rumis/storage/scache"
	"github.com/rumis/storage/srepo"
)

type Zero interface {
	Zero() bool
}

// NewGenericCacheReader 通用缓存读取
func NewGenericCacheReader(prefix string) scache.RedisKeyValueObjectReader {
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

// NewGenericCacheWriter 通用缓存写入
func NewGenericCacheWriter(prefix string) func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
	w := scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()), scache.WithPrefix(prefix))
	return func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
		err := w(ctx, kv, expire)
		return err
	}
}

// NewGenericRepoReader 通用数据库读取
func NewGenericRepoReader(tableName string, columns []string) func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
	r := srepo.NewSealMysqlOneReader(srepo.WithName(tableName), srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns))
	return func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
		err := r(ctx, out, hand...)
		return err
	}
}

// NewKeyValueClauseHandler 查询匹配条件
func NewEqClauseHandler(key string, val interface{}) srepo.ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.Eq(key, val))
	}
}

// NewGenericLocker 通用锁
func NewGenericLocker(biz string) locker.Locker {
	return scache.DefaultRedisLocker(scache.DefaultClient(), biz)
}

// NewGenericReader 通用缓存-库数据读取器
func NewGenericReader(prefix string, tablename string, columns []string, biz string, opts ...srepo.ClauseHandler) func(ctx context.Context, params interface{}, expire time.Duration, out interface{}) error {
	cacheReader := NewGenericCacheReader(prefix)
	cacheWriter := NewGenericCacheWriter(prefix)
	repoReader := NewGenericRepoReader(tablename, columns)
	locker := NewGenericLocker(biz)
	return func(ctx context.Context, params interface{}, expire time.Duration, out interface{}) error {
		zero, ok := out.(Zero)
		if !ok {
			return errors.New("out must implements Zero")
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
