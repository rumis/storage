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

// SQL子句处理方法
type ClauseHandler func(interface{})

// 数据插入
type RepoInserter func(context.Context, interface{}) (int64, error)

// 数据更新
type RepoUpdater func(context.Context, interface{}, ...ClauseHandler) (int64, error)

// 数据读取
type RepoReader func(context.Context, interface{}, ...ClauseHandler) error
