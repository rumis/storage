package srepo

import (
	"context"

	"github.com/rumis/seal"
)

// NewSealMysqlOneReader 创建新的Seal数据写入对象
func NewSealMysqlOneReader(hands ...RepoSealOptionHandler) RepoReader {
	// 默认配置
	opts := DefaultRepoSealOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	// 优先TX
	if sealTx, ok := opts.TX.(*seal.Tx); ok {
		return func(ctx context.Context, data interface{}, handler ...ClauseHandler) error {
			q := sealTx.Select(opts.Columns...).From(opts.Name)
			for _, v := range handler {
				v(q)
			}
			err := q.Limit(1).Query(ctx).OneStruct(&data)
			return err
		}
	}
	// DB逻辑
	if sealDb, ok := opts.DB.(seal.DB); ok {
		return func(ctx context.Context, data interface{}, handler ...ClauseHandler) error {
			q := sealDb.Select(opts.Columns...).From(opts.Name)
			for _, v := range handler {
				v(q)
			}
			err := q.Limit(1).Query(ctx).OneStruct(&data)
			return err
		}
	}
	// error
	return func(ctx context.Context, data interface{}, handler ...ClauseHandler) error {
		return ErrBothDbAndTxNil
	}
}
