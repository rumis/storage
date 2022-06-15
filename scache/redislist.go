package scache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage"
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
		if opts.Client == nil {
			return ErrClientNil
		}
		switch val := params.(type) {
		case string:
			_, err := opts.Client.RPush(ctx, opts.Prefix, val).Result()
			if err != nil {
				return err
			}
		case []string:
			ivals := make([]interface{}, 0, len(val))
			for _, v := range val {
				ivals = append(ivals, v)
			}
			_, err := opts.Client.RPush(ctx, opts.Prefix, ivals...).Result()
			if err != nil {
				return err
			}
			return nil
		case Pair:
			_, err := opts.Client.RPush(ctx, opts.Prefix+val.Key, val.Value).Result()
			if err != nil {
				return err
			}
		case []Pair:
			for _, v := range val {
				_, err := opts.Client.RPush(ctx, opts.Prefix+v.Key, v.Value).Result()
				if err != nil {
					return err
				}
			}
			return nil
		case storage.ForEach:
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
				_, err = opts.Client.RPush(ctx, key, string(val)).Result()
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
			key, err := opts.KeyFn(val)
			if err != nil {
				return err
			}
			buf, err := ujson.Marshal(val)
			if err != nil {
				return err
			}
			_, err = opts.Client.RPush(ctx, key, string(buf)).Result()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// NewRedisListReader 创建新的Redis队列读取对象
func NewRedisListReader(hands ...RedisOptionHandler) RedisListReader {
	// 默认配置
	opts := DefaultRedisOptions()
	// 自定义配置设置
	for _, hand := range hands {
		hand(&opts)
	}
	return func(ctx context.Context) (interface{}, error) {
		if opts.Client == nil {
			return nil, ErrClientNil
		}
		elem, err := opts.Client.LPop(ctx, opts.Prefix).Result()
		if err == redis.Nil {
			return "", nil
		}
		if err != nil {
			return "", err
		}
		return elem, nil
	}
}
