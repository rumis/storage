package scache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/pkg/ujson"
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
		case string:
			cmd := opts.Client.RPush(ctx, opts.Prefix, val)
			_, err := cmd.Result()
			if opts.ExecLogFn != nil {
				opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
			}
			if err != nil {
				return err
			}
		case []string:
			ivals := make([]interface{}, 0, len(val))
			for _, v := range val {
				ivals = append(ivals, v)
			}
			cmd := opts.Client.RPush(ctx, opts.Prefix, ivals...)
			_, err := cmd.Result()
			if opts.ExecLogFn != nil {
				opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
			}
			if err != nil {
				return err
			}
		case Pair:
			cmd := opts.Client.RPush(ctx, opts.Prefix+val.Key, val.Value)
			_, err := cmd.Result()
			if opts.ExecLogFn != nil {
				opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
			}
			if err != nil {
				return err
			}
		case []Pair:
			for _, v := range val {
				startTime = time.Now()
				cmd := opts.Client.RPush(ctx, opts.Prefix+v.Key, v.Value)
				_, err := cmd.Result()
				if opts.ExecLogFn != nil {
					opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
				}
				if err != nil {
					return err
				}
			}
		case meta.ForEach:
			if opts.KeyFn == nil {
				return ErrKeyFnNil
			}
			err := val.ForEach(func(item interface{}) error {
				key, err := opts.KeyFn(item)
				if err != nil {
					return err
				}
				val, err := ujson.Marshal(item)
				if err != nil {
					return err
				}
				startTime = time.Now()
				cmd := opts.Client.RPush(ctx, key, string(val))
				_, err = cmd.Result()
				if opts.ExecLogFn != nil {
					opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
				}
				return err
			})
			if err != nil {
				return err
			}
		default:
			if opts.KeyFn == nil {
				return ErrKeyFnNil
			}
			key, err := opts.KeyFn(val)
			if err != nil {
				return err
			}
			buf, err := ujson.Marshal(val)
			if err != nil {
				return err
			}
			cmd := opts.Client.RPush(ctx, key, string(buf))
			_, err = cmd.Result()
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
func NewRedisListStringReader(hands ...RedisOptionHandler) RedisListStringReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context) (string, error) {
		if opts.Client == nil {
			return "", ErrClientNil
		}
		startTime := time.Now()
		// Exec
		cmd := opts.Client.LPop(ctx, opts.Prefix)
		elem, err := cmd.Result()
		// Log
		if opts.ExecLogFn != nil {
			opts.ExecLogFn(ctx, time.Since(startTime), cmd.String(), err)
		}
		if err == redis.Nil {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		return elem, nil
	}
}

// NewRedisListObjectReader 创建新的Redis List对象读取，读取器返回值为对象
func NewRedisListObjectReader(hands ...RedisOptionHandler) RedisListObjectReader {
	strReader := NewRedisListStringReader(hands...)
	return func(ctx context.Context, data interface{}) error {
		elem, err := strReader(ctx)
		if err != nil {
			return err
		}
		if elem == "" {
			return nil
		}
		err = ujson.Unmarshal([]byte(elem), data)
		if err != nil {
			return err
		}
		return nil
	}
}
