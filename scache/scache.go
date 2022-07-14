package scache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// Pair 键值对
type Pair struct {
	Key   string
	Value string
}

// Redis K-V类型读取
type RedisKeyValueObjectReader func(ctx context.Context, params interface{}, out interface{}) error

// Redis K-V类型读取
type RedisKeyValueStringReader func(ctx context.Context, params interface{}) (interface{}, error)

// Redis K-V类型写入
type RedisKeyValueWriter func(ctx context.Context, params interface{}, expire time.Duration) error

// Redis K-V类型删除
type RedisKeyValueDeleter func(ctx context.Context, param interface{}) error

// RedisKeyGenerator Redis Key生成
type RedisKeyGenerator func(param interface{}) (string, error)

// Redis List类型写入
type RedisListWriter func(context.Context, interface{}) error

// RedisListStringReader Redis List类型读取，返回结果为字符串
type RedisListStringReader func(context.Context) (string, error)

// RedisListObjectReader Redis List类型读取，返回结果为对象
type RedisListObjectReader func(context.Context, interface{}) error

// 选项
type RedisOptions struct {
	KeyFn  RedisKeyGenerator
	Prefix string
	Client *redis.Client
}

// RedisOptionHandler 配置选项
type RedisOptionHandler func(*RedisOptions)

// 创建默认的Redis配置
func DefaultRedisOptions() RedisOptions {
	return RedisOptions{}
}

// WithKeyFn 配置KEY生成方法
func WithKeyFn(fn RedisKeyGenerator) RedisOptionHandler {
	return func(opts *RedisOptions) {
		opts.KeyFn = fn
	}
}

// WithClient 配置客户端实例
func WithClient(client *redis.Client) RedisOptionHandler {
	return func(opts *RedisOptions) {
		opts.Client = client
	}
}

// WithPrefix 配置KEY前缀
func WithPrefix(pre string) RedisOptionHandler {
	return func(opts *RedisOptions) {
		opts.Prefix = pre
	}
}

// Redis客户端空
var ErrClientNil error = errors.New("redis client is nil")
var ErrKeyFnNil error = errors.New("key generater is nil")
var ErrPrefixNil error = errors.New("key prefix is nil")
var ErrKeyGenerate error = errors.New("key generate error")
var ErrKeyFormat error = errors.New("key format error")
