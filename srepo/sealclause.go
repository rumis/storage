package srepo

import (
	"github.com/rumis/seal"
	"github.com/rumis/seal/query"
)

// NewSealEqClause 查询匹配条件
func NewSealEqClause(key string, val interface{}) ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.Eq(key, val))
	}
}
