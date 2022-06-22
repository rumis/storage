package srepo

import (
	"database/sql"
	"fmt"

	"github.com/rumis/seal"
	"github.com/rumis/seal/builder"
	"github.com/rumis/seal/options"
)

var sealOpts []options.SealOptionsFunc

// InitSealOptions 初始化Seal相关选项
func InitSealOptions(fns ...options.SealOptionsFunc) {
	sealOpts = fns
}

// NewSealMysqlDB 创建Seal.Mysql对象
func NewSealMysqlDB(db *sql.DB) seal.DB {
	sealDb, err := seal.OpenWithDB(db, builder.NewMysqlBuilder())
	if err != nil {
		fmt.Println(err)
	}
	return sealDb
}

// type SelectClauseHandler func(q *query.SelectQuery)

// type UpdateClauseHandler func(q *query.UpdateQuery)
