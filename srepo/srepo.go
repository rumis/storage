package srepo

import (
	"context"
	"errors"
)

// 错误定义
var ErrBothDbAndTxNil error = errors.New("both db and tx is nil")

// 选项
type RepoOptions struct {
	TX      interface{}
	DB      interface{}
	Name    string
	Columns []string
}

// RepoOptionHandler 数据库配置选项
type RepoOptionHandler func(*RepoOptions)

// 创建默认的Repo配置
func DefaultRepoOptions() RepoOptions {
	return RepoOptions{}
}

// WithName 表名
func WithName(name string) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.Name = name
	}
}

// WithDB 数据库实例
func WithDB(db interface{}) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.DB = db
	}
}

// WithDB 数据库实例
func WithTX(tx interface{}) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.TX = tx
	}
}

// WithColumns 配置表字段
func WithColumns(columns []string) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.Columns = columns
	}
}

// ClauseHandler SQL子句处理方法
// @params query 查询器对象或者TX、DB等
type ClauseHandler func(query interface{})

// RepoInserter 数据插入
// @params data 需要插入的数据，支持单个数据或者数组
type RepoInserter func(ctx context.Context, data interface{}) (int64, error)

// RepoUpdater 数据更新
// @params data 需要更新的数据，支持map和struct
// @parama where 更新数据的条件
// @return 最后一个自增ID的值
type RepoUpdater func(ctx context.Context, data interface{}, where ...ClauseHandler) (int64, error)

// RepoReader 数据读取
// @params data 承载数据的指针
// @params where 查询字句
type RepoReader func(ctx context.Context, data interface{}, where ...ClauseHandler) error
