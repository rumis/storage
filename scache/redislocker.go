package scache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/locker"
)

var defaultLockerPrefix = "tal_jiaoyan_storage_locker_"

// DefaultRedisLocker 创建基于Redis的分布式锁
func DefaultRedisLocker(client *redis.Client) locker.Locker {
	return locker.NewLocker(
		locker.WithLockerReader(RedisLockerReader(NewRedisKeyValueReader(WithClient(client), WithPrefix(defaultLockerPrefix)))),
		locker.WithLockerWriter(RedisLockerWriter(NewRedisKeyValueWriter(WithClient(client), WithPrefix(defaultLockerPrefix)), locker.DefaultExpire)),
	)
}

// RedisLockerWriter
func RedisLockerWriter(w RedisKeyValueWriter, expire time.Duration) locker.LockerWriter {
	return func(ctx context.Context, key string) error {
		err := w(ctx, Pair{Key: key, Value: key}, expire)
		return err
	}
}

// RedisLockerReader
func RedisLockerReader(r RedisKeyValueReader) locker.LockerReader {
	return func(ctx context.Context, key string) (string, error) {
		val, err := r(ctx, key)
		if err != nil {
			return "", err
		}
		sval, ok := val.(string)
		if !ok {
			return "", errors.New("locker value type error")
		}
		return sval, nil
	}
}
