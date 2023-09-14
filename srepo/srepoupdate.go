package srepo

import (
	"context"

	"github.com/rumis/seal"
	"github.com/rumis/storage/v2/meta"
)

// NewSealMysqlMultiReader 创建新的Seal数据读取对象，返回值多行
func NewSealMysqlUpdater(hands ...RepoSealOptionHandler) meta.RepoUpdater {
	// 默认配置
	opts := DefaultRepoSealOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, param interface{}, handler meta.QueryExprHandler) (int64, error) {
			var affectCnt int64
			q := sealTx.Update(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Value(param).Exec(ctx, &affectCnt)
			return affectCnt, err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, param interface{}, handler meta.QueryExprHandler) (int64, error) {
			var affectCnt int64
			q := sealDb.Update(opts.Name)
			exprs := handler(ctx)
			for _, expr := range exprs {
				q.Where(expr)
			}
			err := q.Value(param).Exec(ctx, &affectCnt)
			return affectCnt, err
		}
	}
	// error
	return func(ctx context.Context, param interface{}, handler meta.QueryExprHandler) (int64, error) {
		return 0, ErrBothDbAndTxNil
	}
}
