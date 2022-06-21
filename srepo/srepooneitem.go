package srepo

import (
	"context"
)

// NewSealMysqlOneReader 创建新的Seal数据写入对象
func NewSealMysqlOneReader(hands ...RepoOptionHandler) RepoReader {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	sealDb := NewSealMysqlDB(opts.DB)
	return func(ctx context.Context, data interface{}, handler ...ClauseHandler) error {
		q := sealDb.Select(opts.Columns...).From(opts.Name)
		for _, v := range handler {
			v(q)
		}
		err := q.Limit(1).Query().OneStruct(&data)
		return err
	}
}
