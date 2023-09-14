package meta

import (
	"context"

	"github.com/rumis/seal/expr"
)

type QueryExprHandler func(context.Context) []expr.Expr

// RepoInserter 数据插入
// @params data 需要插入的数据，支持单个数据或者数组
type RepoInserter func(ctx context.Context, data interface{}) (int64, error)

// RepoUpdater 数据更新
// @params data 需要更新的数据，支持map和struct
// @parama where 更新数据的条件
// @return 最后一个自增ID的值
type RepoUpdater func(ctx context.Context, data interface{}, where QueryExprHandler) (int64, error)

// RepoReader 数据读取
// @params data 承载数据的指针
// @params where 查询字句
type RepoReader func(ctx context.Context, data interface{}, where QueryExprHandler) error
