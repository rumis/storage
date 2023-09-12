package scache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/locker"
)

var defaultLockerPrefix = "tal_jiaoyan_storage_locker_"

// DefaultRedisLocker 创建基于Redis的分布式锁
func DefaultRedisLocker(client *redis.Client, biz string) locker.Locker {
	return locker.NewLocker(
		locker.WithLockerAdder(RedisLockerAdder(NewRedisKeyValueSetNX(WithClient(client)), locker.DefaultExpire)),
		locker.WithLockerDeleter(RedisLockerDeleter(NewRedisKeyValueDeleter(WithClient(client)))),
	)
}

// RedisLockerWriter
func RedisLockerAdder(a RedisKeyValueSetNX, expire time.Duration) locker.LockerAdder {
	return func(ctx context.Context, key string) bool {
		err := a(ctx, StringKey(key), expire)
		return err
	}
}

// RedisLockerDeleter
func RedisLockerDeleter(d RedisKeyValueDeleter) locker.LockerDeleter {
	return func(ctx context.Context, key string) error {
		err := d(ctx, StringKey(key))
		return err
	}
}
