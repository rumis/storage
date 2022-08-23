package scache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/meta"
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
//
// 参数param支持以下4种类型:
//
// 	string: 需要和prefix参数配合
// 	[]string: 需要和prefix参数配合
// 	实现接口ForEach：需要和KeyFn参数配合
// 	其他值：需要和KeyFn参数配合
type RedisKeyValueDeleter func(ctx context.Context, param interface{}) error

// RedisKeyGenerator Redis Key生成
type RedisKeyGenerator func(param interface{}) (string, error)

// Redis List类型写入

// 参数param支持如下6种类型：

// 	string：需要和prefix配合，自动写入key为【prefix】的List中
// 	[]string: 需要和prefix配合，自动写入key为【prefix】的List中
// 	Pair：需要和prefix配合，写入key为【prefix+p.Key】的List中
// 	[]Pair：需要和prefix配合，写入key为【prefix+p.Key】的List中
// 	实现接口ForEach：需要和KeyFn参数配合,每个元素通过KeyFn计算key，然后写入对应的List中，值经过json序列化
// 	其他值：需要和KeyFn参数配合，通过KeyFn计算key，然后写入【key】的List中
type RedisListWriter func(context.Context, interface{}) error

// RedisListStringReader Redis List类型读取，每次读取一个值，返回结果为字符串
type RedisListStringReader func(context.Context) (string, error)

// RedisListObjectReader Redis List类型读取，每次读取一个值，返回结果为对象
type RedisListObjectReader func(context.Context, interface{}) error

// ExecLogError 记录调用日志
func ExecLogError(ctx context.Context, fn meta.RedisExecLogFunc, stime time.Time, args interface{}, e error) error {
	if fn != nil {
		fn(ctx, time.Since(stime), args, e)
	}
	return e
}

// 选项
type RedisOptions struct {
	KeyFn     RedisKeyGenerator
	Prefix    string
	Client    *redis.Client
	ExecLogFn meta.RedisExecLogFunc
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

// WithPrefix 配置KEY前缀
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
