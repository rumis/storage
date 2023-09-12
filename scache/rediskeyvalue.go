package scache

import (
	"context"
	"strings"
	"time"

	"github.com/rumis/storage/v2/meta"
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
		case meta.ForEach:
			err := vals.ForEach(func(item interface{}) error {
				keyIT, ok := item.(meta.Key)
				if !ok {
					return meta.EI_KeyNotImplement
				}
				key := keyIT.Key()
				valIT, ok := item.(meta.String)
				if !ok {
					return meta.EI_StringNotImplement
				}
				val := valIT.String()
				cmd := opts.Client.Set(ctx, key, val, expiration)
				err := cmd.Err()
				if err != nil {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
				}
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
			})
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
		default:
			keyIT, ok := params.(meta.Key)
			if !ok {
				return meta.EI_KeyNotImplement
			}
			key := keyIT.Key()
			valIT, ok := params.(meta.String)
			if !ok {
				return meta.EI_StringNotImplement
			}
			val := valIT.String()
			cmd := opts.Client.Set(ctx, key, string(val), expiration)
			err := cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		}
		return ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
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
	return func(ctx context.Context, params interface{}, data interface{}) error {
		startTime := time.Now()
		if opts.Client == nil {
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, ErrClientNil)
		}
		switch keys := params.(type) {
		case meta.ForEach:
			allRes := make([]string, 0)
			err := keys.ForEach(func(item interface{}) error {
				keyIT, ok := item.(meta.Key)
				if !ok {
					return ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_KeyNotImplement)
				}
				key := keyIT.Key()
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
			valueIT, ok := params.(meta.Value)
			if !ok {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_ValueNotImplement)
			}
			err = valueIT.Value(totalRes)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
			}
			return ExecLogError(ctx, opts.ExecLogFn, startTime, params, nil)
		default:
			keyIT, ok := params.(meta.Key)
			if !ok {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_KeyNotImplement)
			}
			key := keyIT.Key()
			cmd := opts.Client.Get(ctx, key)
			res, err := cmd.Result()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			valueIT, ok := params.(meta.Value)
			if !ok {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_ValueNotImplement)
			}
			err = valueIT.Value(res)
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, params, err)
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
		case meta.ForEach:
			keys := make([]string, 0)
			err := vals.ForEach(func(item interface{}) error {
				keyIT, ok := item.(meta.Key)
				if !ok {
					return meta.EI_KeyNotImplement
				}
				key := keyIT.Key()
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
			keyIT, ok := params.(meta.Key)
			if !ok {
				return meta.EI_KeyNotImplement
			}
			key := keyIT.Key()
			cmd := opts.Client.Del(ctx, key)
			err := cmd.Err()
			if err != nil {
				return ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			}
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		}
		return nil
	}
}

// NewRedisKeyValueSetNX 创建SetNX方法
func NewRedisKeyValueSetNX(hands ...RedisOptionHandler) RedisKeyValueSetNX {
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
		keyIT, ok := params.(meta.Key)
		if !ok {
			ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_KeyNotImplement)
			return false
		}
		key := keyIT.Key()
		valueIT, ok := params.(meta.String)
		if !ok {
			ExecLogError(ctx, opts.ExecLogFn, startTime, params, meta.EI_StringNotImplement)
		}
		val := valueIT.String()
		cmd := opts.Client.SetNX(ctx, key, val, expiration)
		err := cmd.Err()
		if err != nil {
			ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), err)
			return false
		}
		ExecLogError(ctx, opts.ExecLogFn, startTime, cmd.String(), nil)
		return cmd.Val()

	}
}
