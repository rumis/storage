package srepo

import (
	"context"

	"github.com/rumis/seal"
	"github.com/rumis/storage/v2/meta"
)

// NewSealMysqlOneReader 创建新的Seal数据写入对象
func NewSealMysqlOneReader(hands ...RepoSealOptionHandler) meta.RepoReader {
	// 默认配置
	opts := DefaultRepoSealOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
			q := sealTx.Select(opts.Columns...).From(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Limit(1).Query(ctx).OneStruct(&data)
			return err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
			q := sealDb.Select(opts.Columns...).From(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Limit(1).Query(ctx).OneStruct(&data)
			return err
		}
	}
	// error
	return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
		return ErrBothDbAndTxNil
	}
}

// NewSealMysqlMultiReader 创建新的Seal数据读取对象，返回值多行
func NewSealMysqlMultiReader(hands ...RepoSealOptionHandler) meta.RepoReader {
	// 默认配置
	opts := DefaultRepoSealOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
			q := sealTx.Select(opts.Columns...).From(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Query(ctx).AllStruct(data)
			return err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
			q := sealDb.Select(opts.Columns...).From(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Query(ctx).AllStruct(data)
			return err
		}
	}
	// error
	return func(ctx context.Context, data interface{}, handler meta.QueryExprHandler) error {
		return ErrBothDbAndTxNil
	}
}
