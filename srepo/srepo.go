package srepo

import (
	"context"
	"errors"

	"github.com/rumis/storage/v2/meta"
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

// RepoGroupReader 多表数据读取
// @params data 承载数据的指针
// @params params 查询条件字段
type RepoGroupReader func(context.Context, interface{}, interface{}) error

// RepoGroupOptions 自定义复杂操作集合
type RepoGroupOptions struct {
	Handler     interface{}
	ExecLogFunc meta.RepoExecLogFunc
}

// RepoGroupOptionHandler 数据库配置选项
type RepoGroupOptionHandler func(*RepoGroupOptions)

// DefaultRepoGroupOptions 创建默认的Repo配置
func DefaultRepoGroupOptions() RepoGroupOptions {
	return RepoGroupOptions{}
}

// WithHandler 处理函数
func WithHandler(h interface{}) RepoGroupOptionHandler {
	return func(opts *RepoGroupOptions) {
		opts.Handler = h
	}
}

// WithExecLogger 日志函数
func WithExecLogger(fn meta.RepoExecLogFunc) RepoGroupOptionHandler {
	return func(opts *RepoGroupOptions) {
		opts.ExecLogFunc = fn
	}
}
