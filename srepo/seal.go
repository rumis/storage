package srepo

import (
	"database/sql"
	"fmt"

	"github.com/rumis/seal"
	"github.com/rumis/seal/builder"
)

// NewSealMysqlDB 创建Seal.Mysql对象
func NewSealMysqlDB(db *sql.DB) seal.DB {
	sealDb, err := seal.OpenWithDB(db, builder.NewMysqlBuilder())
	if err != nil {
		fmt.Println(err)
	}
	return sealDb
}

type ClauseHandler func(interface{})

// type SelectClauseHandler func(q *query.SelectQuery)

// type UpdateClauseHandler func(q *query.UpdateQuery)
