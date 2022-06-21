package srepo

import "context"

// NewSealMysqlMultiReader 创建新的Seal数据读取对象，返回值多行
func NewSealMysqlUpdater(hands ...RepoOptionHandler) RepoUpdater {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	sealDb := NewSealMysqlDB(opts.DB)
	return func(ctx context.Context, param interface{}, handler ...ClauseHandler) (int64, error) {
		var affectCnt int64
		q := sealDb.Update(opts.Name)
		for _, v := range handler {
			v(q)
		}
		err := q.Value(param).Exec(&affectCnt)
		return affectCnt, err
	}
}
