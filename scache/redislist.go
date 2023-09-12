package scache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/meta"
)

// NewRedisListWriter 创建新的Redis队列写入对象
func NewRedisListWriter(hands ...RedisOptionHandler) RedisListWriter {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, params interface{}) error {
		startTime := time.Now()
		if opts.Client == nil {
			return ErrClientNil
		}
		switch val := params.(type) {
		case meta.ForEach:
			err := val.ForEach(func(item interface{}) error {
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
				cmd := opts.Client.RPush(ctx, key, val)
				_, err := cmd.Result()
				if opts.ExecLogFn != nil {
					opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
				}
				return err
			})
			if err != nil {
				return err
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
			val1 := valIT.String()
			cmd := opts.Client.RPush(ctx, key, val1)
			_, err := cmd.Result()
			if opts.ExecLogFn != nil {
				opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// NewRedisListReader 创建新的Redis队列读取，读取器返回值为字符串
func NewRedisListReader(hands ...RedisOptionHandler) RedisListReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context, out interface{}) error {
		if opts.Client == nil {
			return ErrClientNil
		}
		startTime := time.Now()
		// 获取缓存key
		keyIT, ok := out.(meta.Key)
		if !ok {
			return meta.EI_KeyNotImplement
		}
		key := keyIT.Key()

		// 读取数据
		cmd := opts.Client.LPop(ctx, key)
		elem, err := cmd.Result()
		if err == redis.Nil {
			return nil
		}
		if err != nil {
			return err
		}
		// 初始化对象内容
		valIT, ok := out.(meta.Value)
		if !ok {
			return meta.EI_ValueNotImplement
		}
		err = valIT.Value(elem)
		if err != nil {
			return err
		}
		// Log
		if opts.ExecLogFn != nil {
			opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
		}
		return nil
	}
}
