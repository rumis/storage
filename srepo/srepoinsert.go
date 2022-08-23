package srepo

import (
	"context"

	"github.com/rumis/seal"
)

// NewSealMysqlInserter 创建新的Seal数据写入对象
func NewSealMysqlInserter(hands ...RepoOptionHandler) RepoInserter {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, params interface{}) (int64, error) {
			var lastId int64
			err := sealTx.Insert(opts.Name).Value(params).Exec(ctx, &lastId)
			return lastId, err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, params interface{}) (int64, error) {
			var lastId int64
			err := sealDb.Insert(opts.Name).Value(params).Exec(ctx, &lastId)
			return lastId, err
		}
	}
	// error
	return func(ctx context.Context, params interface{}) (int64, error) {
		return 0, ErrBothDbAndTxNil
	}
}

// NewSealMysqlMultiInserter 创建新的Seal数据写入对象-一次写入多次数据
func NewSealMysqlMultiInserter(hands ...RepoOptionHandler) RepoInserter {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, params interface{}) (int64, error) {
			var lastId int64
			err := sealTx.Insert(opts.Name).Values(params).Exec(ctx, &lastId)
			return lastId, err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, params interface{}) (int64, error) {
			var lastId int64
			err := sealDb.Insert(opts.Name).Values(params).Exec(ctx, &lastId)
			return lastId, err
		}
	}
	// error
	return func(ctx context.Context, params interface{}) (int64, error) {
		return 0, ErrBothDbAndTxNil
	}
}
