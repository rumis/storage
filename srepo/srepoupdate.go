package srepo

import (
	"context"

	"github.com/rumis/seal"
)

// NewSealMysqlMultiReader 创建新的Seal数据读取对象，返回值多行
func NewSealMysqlUpdater(hands ...RepoOptionHandler) RepoUpdater {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, param interface{}, handler ...ClauseHandler) (int64, error) {
			var affectCnt int64
			q := sealTx.Update(opts.Name)
			for _, v := range handler {
				v(q)
			}
			err := q.Value(param).Exec(ctx, &affectCnt)
			return affectCnt, err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, param interface{}, handler ...ClauseHandler) (int64, error) {
			var affectCnt int64
			q := sealDb.Update(opts.Name)
			for _, v := range handler {
				v(q)
			}
			err := q.Value(param).Exec(ctx, &affectCnt)
			return affectCnt, err
		}
	}
	// error
	return func(ctx context.Context, param interface{}, handler ...ClauseHandler) (int64, error) {
		return 0, ErrBothDbAndTxNil
	}
}
