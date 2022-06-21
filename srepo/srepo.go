package srepo

import (
	"context"
	"database/sql"
)

// 选项
type RepoOptions struct {
	DB      *sql.DB
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
func WithDB(db *sql.DB) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.DB = db
	}
}

// WithColumns 配置表字段
func WithColumns(columns []string) RepoOptionHandler {
	return func(opts *RepoOptions) {
		opts.Columns = columns
	}
}

// 数据插入
type RepoInster func(context.Context, interface{}) (int64, error)

// 数据更新
type RepoUpdater func(context.Context, interface{}, ...ClauseHandler) (int64, error)

// 数据读取
type RepoReader func(context.Context, interface{}, ...ClauseHandler) error
