package srepo

import (
	"context"
	"errors"
	"time"
)

// NewSealMysqlOneReader 创建新的Seal数据写入对象
func NewMysqlGroupReader(hands ...RepoGroupOptionHandler) RepoGroupReader {
	// 默认配置
	opts := DefaultRepoGroupOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	return func(ctx context.Context, data interface{}, params interface{}) error {
		startTime := time.Now()
		switch fn := opts.Handler.(type) {
		case RepoInserter:
			res, ok := data.(*int64)
			if !ok {
				err := errors.New("type of params data must be int64 pointer")
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
				}
				return err
			}
			lastId, err := fn(ctx, params)
			if err != nil {
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
				}
				return err
			}
			*res = lastId // 赋值
		case RepoUpdater:
			var where []ClauseHandler
			var ok bool
			if params != nil {
				where, ok = params.([]ClauseHandler)
				if !ok {
					err := errors.New("type of params must be ClauseHandler slice")
					if opts.ExecLogFunc != nil {
						opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
					}
					return err
				}
			}
			affectCount, err := fn(ctx, data, where...)
			if err != nil {
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
				}
				return err
			}
			// 如果影响行数为零，通过error返回
			if affectCount == 0 {
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, ErrUpdateAffectZeroRows)
				}
				return ErrUpdateAffectZeroRows
			}
		case RepoReader:
			var where []ClauseHandler
			var ok bool
			if params != nil {
				where, ok = params.([]ClauseHandler)
				if !ok {
					err := errors.New("type of params must be ClauseHandler slice")
					if opts.ExecLogFunc != nil {
						opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
					}
					return err
				}
			}
			err := fn(ctx, data, where...)
			if err != nil {
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
				}
				return err
			}
		case RepoGroupReader:
			err := fn(ctx, data, params)
			if err != nil {
				if opts.ExecLogFunc != nil {
					opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
				}
				return err
			}
		default:
			err := errors.New("unsupport handler")
			if opts.ExecLogFunc != nil {
				opts.ExecLogFunc(ctx, time.Since(startTime), params, err)
			}
			return err
		}
		if opts.ExecLogFunc != nil {
			opts.ExecLogFunc(ctx, time.Since(startTime), params, nil)
		}
		return nil
	}
}
