package srepo

import "context"

// NewSealMysqlInserter 创建新的Seal数据写入对象
func NewSealMysqlInserter(hands ...RepoOptionHandler) RepoInster {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	sealDb := NewSealMysqlDB(opts.DB)
	return func(ctx context.Context, params interface{}) (int64, error) {
		var lastId int64
		err := sealDb.Insert(opts.Name).Value(params).Exec(&lastId)
		return lastId, err
	}
}

// NewSealMysqlMultiInserter 创建新的Seal数据写入对象-一次写入多次数据
func NewSealMysqlMultiInserter(hands ...RepoOptionHandler) RepoInster {
	// 默认配置
	opts := DefaultRepoOptions()
	// 自定义配置设置
	for _, fn := range hands {
		fn(&opts)
	}
	sealDb := NewSealMysqlDB(opts.DB)
	return func(ctx context.Context, params interface{}) (int64, error) {
		var lastId int64
		err := sealDb.Insert(opts.Name).Values(params).Exec(&lastId)
		return lastId, err
	}
}
