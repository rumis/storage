package scache

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/pkg/ujson"
)

// NewRedisKeyValueWriter 创建新的缓存写入
func NewRedisKeyValueWriter(hands ...RedisOptionHandler) RedisKeyValueWriter {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}, expiration time.Duration) error {
		startTime := time.Now()
		if opts.Client == nil {
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
		}
		switch vals := params.(type) {
		case Pair:
			if opts.Prefix == "" {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			cmd := opts.Client.Set(ctx, opts.Prefix+vals.Key, vals.Value, expiration)
			err := cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		case []Pair:
			if opts.Prefix == "" {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			for _, val := range vals {
				startTime = time.Now()
				cmd := opts.Client.Set(ctx, opts.Prefix+val.Key, val.Value, expiration)
				err := cmd.Err()
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				}
				ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			}
		case meta.ForEach:
			if opts.KeyFn == nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			err := vals.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				}
				val, err := ujson.Marshal(item)
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				}
				cmd := opts.Client.Set(ctx, key, string(val), expiration)
				err = cmd.Err()
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				}
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			})
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
		default:
			if opts.KeyFn == nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			key, err := opts.KeyFn(vals)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			val, err := ujson.Marshal(vals)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			cmd := opts.Client.Set(ctx, key, string(val), expiration)
			err = cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		}
		return ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
	}
}

// NewRedisKeyValueReader 自定义Redis读取
func NewRedisKeyValueStringReader(hands ...RedisOptionHandler) RedisKeyValueStringReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}) (interface{}, error) {
		startTime := time.Now()
		if opts.Client == nil {
			return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
		}
		switch keys := params.(type) {
		case string:
			if opts.Prefix == "" {
				return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			cmd := opts.Client.Get(ctx, opts.Prefix+keys)
			res, err := cmd.Result()
			if err != nil {
				return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			return res, ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		case []string:
			if opts.Prefix == "" {
				return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			allRes := make([]string, 0, len(keys))
			for _, key := range keys {
				cmd := opts.Client.Get(ctx, opts.Prefix+key)
				res, err := cmd.Result()
				if err == redis.Nil {
					allRes = append(allRes, res)
					continue
				}
				if err != nil {
					return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				}
				allRes = append(allRes, res)
				ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			}
			return allRes, ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
		default:
			return nil, ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFormat)
		}
	}
}

// NewRedisKeyValueReader 自定义Redis读取
func NewRedisKeyValueObjectReader(hands ...RedisOptionHandler) RedisKeyValueObjectReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}, data interface{}) error {
		startTime := time.Now()
		if opts.Client == nil {
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
		}
		switch keys := params.(type) {
		case meta.ForEach:
			if opts.KeyFn == nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			allRes := make([]string, 0)
			err := keys.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				}
				startTime = time.Now()
				cmd := opts.Client.Get(ctx, key)
				res, err := cmd.Result()
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				}
				allRes = append(allRes, res)
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			})
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			// 拼接为数组
			totalRes := "[" + strings.Join(allRes, ",") + "]"
			err = ujson.Unmarshal([]byte(totalRes), data)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
		default:
			if opts.KeyFn == nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			key, err := opts.KeyFn(keys)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			cmd := opts.Client.Get(ctx, key)
			res, err := cmd.Result()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			err = ujson.Unmarshal([]byte(res), data)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
		}
	}
}

// NewRedisKeyValueDeleter 缓存删除
func NewRedisKeyValueDeleter(hands ...RedisOptionHandler) RedisKeyValueDeleter {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}) error {
		startTime := time.Now()
		if opts.Client == nil {
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
		}
		switch vals := params.(type) {
		case string:
			if opts.Prefix == "" {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			cmd := opts.Client.Del(ctx, opts.Prefix+vals)
			err := cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		case []string:
			if opts.Prefix == "" {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
			}
			for i := range vals {
				vals[i] = opts.Prefix + vals[i]
			}
			cmd := opts.Client.Del(ctx, vals...)
			err := cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		case meta.ForEach:
			if opts.KeyFn == nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			keys := make([]string, 0)
			err := vals.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return err
				}
				keys = append(keys, key)
				return nil
			})
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			cmd := opts.Client.Del(ctx, keys...)
			err = cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		default:
			if opts.KeyFn == nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
			}
			key, err := opts.KeyFn(vals)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			cmd := opts.Client.Del(ctx, key)
			err = cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		}
		return nil
	}
}

func NewReaderSetNX(hands ...RedisOptionHandler) RedisKeyValueNX {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}, expiration time.Duration) bool {
		startTime := time.Now()
		if opts.Client == nil {
			ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
			return false
		}
		switch vals := params.(type) {
		case Pair:
			if opts.Prefix == "" {
				ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrPrefixNil)
				return false
			}
			cmd := opts.Client.SetNX(ctx, opts.Prefix+vals.Key, vals.Value, expiration)
			err := cmd.Err()
			if err != nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				return false
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			return cmd.Val()
		default:
			if opts.KeyFn == nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrKeyFnNil)
				return false
			}
			key, err := opts.KeyFn(vals)
			if err != nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				return false
			}
			val, err := ujson.Marshal(vals)
			if err != nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
				return false
			}
			cmd := opts.Client.SetNX(ctx, key, string(val), expiration)
			err = cmd.Err()
			if err != nil {
				ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				return false
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			return cmd.Val()
		}
	}
}
