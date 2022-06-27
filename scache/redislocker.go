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
func DefaultRedisLocker(client *redis.Client, biz string) locker.Locker {
	prefix := defaultLockerPrefix + biz + "_"
	return locker.NewLocker(
		locker.WithLockerReader(RedisLockerReader(NewRedisKeyValueStringReader(WithClient(client), WithPrefix(prefix)))),
		locker.WithLockerWriter(RedisLockerWriter(NewRedisKeyValueWriter(WithClient(client), WithPrefix(prefix)), locker.DefaultExpire)),
		locker.WithLockerDeleter(RedisLockerDeleter(NewRedisKeyValueDeleter(WithClient(client), WithPrefix(prefix)))),
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
func RedisLockerReader(r RedisKeyValueStringReader) locker.LockerReader {
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

// RedisLockerDeleter
func RedisLockerDeleter(d RedisKeyValueDeleter) locker.LockerDeleter {
	return func(ctx context.Context, key string) error {
		err := d(ctx, key)
		return err
	}
}
