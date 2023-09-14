package srepo

import (
	"errors"
)

// 错误定义
var ErrBothDbAndTxNil error = errors.New("both db and tx is nil")
var ErrUpdateAffectZeroRows error = errors.New("update clauses affect zero rows")

// 选项
type RepoSealOptions struct {
	TX      interface{}
	DB      interface{}
	Name    string
	Columns []string
}

// RepoSealOptionHandler Seal数据库配置选项
type RepoSealOptionHandler func(*RepoSealOptions)

// 创建默认的Seal Repo配置
func DefaultRepoSealOptions() RepoSealOptions {
	return RepoSealOptions{}
}

// WithName 表名
func WithName(name string) RepoSealOptionHandler {
	return func(opts *RepoSealOptions) {
		opts.Name = name
	}
}

// WithDB 数据库实例
func WithDB(db interface{}) RepoSealOptionHandler {
	return func(opts *RepoSealOptions) {
		opts.DB = db
	}
}

// WithDB 数据库实例
func WithTX(tx interface{}) RepoSealOptionHandler {
	return func(opts *RepoSealOptions) {
		opts.TX = tx
	}
}

// WithColumns 配置表字段
func WithColumns(columns []string) RepoSealOptionHandler {
	return func(opts *RepoSealOptions) {
		opts.Columns = columns
	}
}
