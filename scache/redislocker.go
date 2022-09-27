package scache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/locker"
)

var defaultLockerPrefix = "tal_jiaoyan_storage_locker_"

// DefaultRedisLocker 创建基于Redis的分布式锁
func DefaultRedisLocker(client *redis.Client, biz string) locker.Locker {
	prefix := defaultLockerPrefix + biz + "_"
	return locker.NewLocker(
		locker.WithLockerAdder(RedisLockerAdder(NewReaderSetNX(WithClient(client), WithPrefix(prefix)), locker.DefaultExpire)),
		locker.WithLockerDeleter(RedisLockerDeleter(NewRedisKeyValueDeleter(WithClient(client), WithPrefix(prefix)))),
	)
}

// RedisLockerWriter
func RedisLockerAdder(a RedisKeyValueNX, expire time.Duration) locker.LockerAdder {
	return func(ctx context.Context, key string) bool {
		err := a(ctx, Pair{Key: key, Value: key}, expire)
		return err
	}
}

// RedisLockerDeleter
func RedisLockerDeleter(d RedisKeyValueDeleter) locker.LockerDeleter {
	return func(ctx context.Context, key string) error {
		err := d(ctx, key)
		return err
	}
}
