package scache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/meta"
)

// StringKey String包装，用于Redis KEY
type StringKey string

// Key 提供key
func (s StringKey) Key() string {
	return string(s)
}

type StringKeySlice []StringKey

func (ss StringKeySlice) ForEach(f meta.Iterator) error {
	for _, v := range ss {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Redis K-V类型读取
type RedisKeyValueReader func(ctx context.Context, param interface{}, out interface{}) error

// Redis K-V类型写入
type RedisKeyValueWriter func(ctx context.Context, param interface{}, expire time.Duration) error

// RedisKeyValueNX Redis SetNX
type RedisKeyValueSetNX func(ctx context.Context, param interface{}, expire time.Duration) bool

// RedisKeyValueSetExp 设置超时时间
type RedisKeyValueSetExp func(ctx context.Context, param interface{}, expire time.Duration) error

// Redis K-V类型删除
type RedisKeyValueDeleter func(ctx context.Context, param interface{}) error

// Redis List类型写入
type RedisListWriter func(ctx context.Context, param interface{}) error

// RedisListReader Redis List类型读取，每次读取一个值
type RedisListReader func(ctx context.Context, out interface{}) error

// ExecLogError 记录调用日志
func ExecLogError(ctx context.Context, fn meta.RedisExecLogFunc, stime time.Time, args interface{}, e error) error {
	if fn != nil {
		fn(ctx, time.Since(stime), args, e)
	}
	return e
}

// 选项
type RedisOptions struct {
	Client    *redis.Client
	ExecLogFn meta.RedisExecLogFunc
}

// RedisOptionHandler 配置选项
type RedisOptionHandler func(*RedisOptions)

// 创建默认的Redis配置
func DefaultRedisOptions() RedisOptions {
	return RedisOptions{}
}

// WithClient 配置客户端实例
func WithClient(client *redis.Client) RedisOptionHandler {
	return func(opts *RedisOptions) {
		opts.Client = client
	}
}

// WithExecLogger 配置日志记录方法
func WithExecLogger(fn meta.RedisExecLogFunc) RedisOptionHandler {
	return func(opts *RedisOptions) {
		opts.ExecLogFn = fn
	}
}

// Redis客户端空
var ErrClientNil error = errors.New("redis client is nil")
var ErrKeyFnNil error = errors.New("key generater is nil")
var ErrPrefixNil error = errors.New("key prefix is nil")
var ErrKeyGenerate error = errors.New("key generate error")
var ErrKeyFormat error = errors.New("key format error")
