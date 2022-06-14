package scache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage"
	"github.com/rumis/storage/pkg/ujson"
)

// KeyValue 键值对
type KeyValue struct {
	Key   string
	Value string
}

// NewRedisKeyValueWriter 创建新的缓存写入
func NewRedisKeyValueWriter(hands ...RedisOptionHandler) RedisKeyValueWriter {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}, expiration time.Duration) error {
		if opts.Client == nil {
			return ErrClientNil
		}
		switch vals := params.(type) {
		case KeyValue:
			if opts.Prefix == "" {
				return ErrPrefixNil
			}
			err := opts.Client.Set(ctx, opts.Prefix+vals.Key, vals.Value, expiration).Err()
			if err != nil {
				return err
			}
		case []KeyValue:
			if opts.Prefix == "" {
				return ErrPrefixNil
			}
			for _, val := range vals {
				err := opts.Client.Set(ctx, opts.Prefix+val.Key, val.Value, expiration).Err()
				if err != nil {
					return err
				}
			}
		case storage.ForEach:
			if opts.KeyFn == nil {
				return ErrKeyFnNil
			}
			err := vals.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return err
				}
				val, err := ujson.Marshal(item)
				if err != nil {
					return err
				}
				err = opts.Client.Set(ctx, key, string(val), expiration).Err()
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		default:
			if opts.KeyFn == nil {
				return ErrKeyFnNil
			}
			key, err := opts.KeyFn(vals)
			if err != nil {
				return err
			}
			val, err := ujson.Marshal(vals)
			if err != nil {
				return err
			}
			err = opts.Client.Set(ctx, key, string(val), expiration).Err()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// NewRedisKeyValueReader 自定义Redis读取
func NewRedisKeyValueReader(hands ...RedisOptionHandler) RedisKeyValueReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}) (interface{}, error) {
		if opts.Client == nil {
			return nil, ErrClientNil
		}
		switch keys := params.(type) {
		case string:
			if opts.Prefix == "" {
				return nil, ErrPrefixNil
			}
			res, err := opts.Client.Get(ctx, opts.Prefix+keys).Result()
			if err != nil {
				return nil, err
			}
			return res, nil
		case []string:
			if opts.Prefix == "" {
				return nil, ErrPrefixNil
			}
			allRes := make([]string, 0, len(keys))
			for _, key := range keys {
				res, err := opts.Client.Get(ctx, opts.Prefix+key).Result()
				if err == redis.Nil {
					allRes = append(allRes, res)
					continue
				}
				if err != nil {
					return nil, err
				}
				allRes = append(allRes, res)
			}
			return allRes, nil
		case storage.ForEach:
			if opts.KeyFn == nil {
				return nil, ErrKeyFnNil
			}
			allRes := make([]string, 0)
			err := keys.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return err
				}
				res, err := opts.Client.Get(ctx, key).Result()
				if err != nil {
					return err
				}
				allRes = append(allRes, res)
				return nil
			})
			if err != nil {
				return allRes, err
			}
			return allRes, nil
		default:
			if opts.KeyFn == nil {
				return nil, ErrKeyFnNil
			}
			key, err := opts.KeyFn(keys)
			if err != nil {
				return nil, err
			}
			res, err := opts.Client.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
}
