package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/pkg/ujson"
	"github.com/rumis/storage/scache"
	"github.com/rumis/storage/srepo"
)

// NewMultiCacheRepoReader 多key读取
// @params keyfn Key生成函数
func NewMultiCacheReader(keyfn scache.RedisKeyGenerator) scache.RedisKeyValueObjectReader {
	mr := scache.NewRedisKeyValueObjectReader(scache.WithClient(scache.DefaultClient()), scache.WithKeyFn(keyfn))
	// @params params 要求实现ForEach接口
	// @params out 要求实现Key、Zero接口
	return func(ctx context.Context, params interface{}, out interface{}) error {
		err := mr(ctx, params, out)
		return err
	}
}

// NewOneCacheWriter 缓存对象写入
func NewMultiCacheWriter(keyfn scache.RedisKeyGenerator) func(ctx context.Context, data interface{}, expire time.Duration) error {
	w := scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()), scache.WithKeyFn(keyfn))
	return func(ctx context.Context, data interface{}, expire time.Duration) error {
		err := w(ctx, data, expire)
		return err
	}
}

// NewMultiRepoReader 通用数据库读取
func NewMultiRepoReader(tableName string, columns []string) func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
	r := srepo.NewSealMysqlMultiReader(srepo.WithName(tableName), srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns))
	return func(ctx context.Context, out interface{}, hand ...srepo.ClauseHandler) error {
		err := r(ctx, out, hand...)
		return err
	}
}

// NewMultiCacheRepoReader1 通用缓存-库数据读取器,多值
// DO NOT USE
// TODO
func NewMultiCacheRepoReader1(keyFn scache.RedisKeyGenerator, tablename string, columns []string, biz string) func(ctx context.Context, params interface{}, expire time.Duration, out interface{}, opts ...srepo.ClauseHandler) error {
	cacheReader := NewMultiCacheReader(keyFn)
	cacheWriter := NewMultiCacheWriter(keyFn)
	repoReader := NewMultiRepoReader(tablename, columns)
	locker := NewDefaultLocker(biz)
	return func(ctx context.Context, params interface{}, expire time.Duration, out interface{}, opts ...srepo.ClauseHandler) error {
		outEach, ok := out.(meta.ForEach)
		if !ok {
			return errors.New("param out must implements ForEach interface")
		}
		paramEach, ok := params.(meta.ForEach)
		if !ok {
			return errors.New("param params must implements ForEach interface")
		}
		// 读取缓存
		err := cacheReader(ctx, params, out)
		if err != nil && err != redis.Nil {
			// 缓存读取错误
			fmt.Println("记录错误")
		}
		// 包含全部数据

		// 检查缓存数据的缺失
		cachedData := make(map[string]interface{})
		outEach.ForEach(func(v interface{}) error {
			zero, ok := v.(meta.Zero)
			if !ok {
				return errors.New("out element must implements Zero interface")
			}
			if zero.Zero() {
				return nil
			}
			key, ok := v.(meta.Key)
			if !ok {
				return errors.New("out element must implements Key interface")
			}
			cachedData[key.Key()] = v
			return nil
		})
		noCachedParam := make([]interface{}, 0)
		paramEach.ForEach(func(v interface{}) error {
			key := fmt.Sprint(v)
			_, ok := cachedData[key]
			if !ok {
				noCachedParam = append(noCachedParam, v)
			}
			return nil
		})
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
		// if zero.Zero() {
		// 	// 写入个空数据
		// 	expire = locker.Expire
		// }
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
