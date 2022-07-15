package srepo

import (
	"github.com/rumis/seal"
	"github.com/rumis/seal/query"
)

// SealQEq 相等
func SealQEq(key string, val interface{}) ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.Eq(key, val))
	}
}

// SealQIn    IN
func SealQIn(key string, val ...interface{}) ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.In(key, val...))
	}
}

// SealQLike 模糊查询
func SealQLike(key string, val string) ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.Like(key, val))
	}
}

// SealQOp 一般操作符 > < >= <= 等
func SealQOp(key string, op string, val interface{}) ClauseHandler {
	return func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			return
		}
		sq.Where(seal.Op(key, op, val))
	}
}

// SealUEq 相等
func SealUEq(key string, val interface{}) ClauseHandler {
	return func(q interface{}) {
		uq, ok := q.(*query.UpdateQuery)
		if !ok {
			return
		}
		uq.Where(seal.Eq(key, val))
	}
}

// SealUIn    IN
func SealUIn(key string, val ...interface{}) ClauseHandler {
	return func(q interface{}) {
		uq, ok := q.(*query.UpdateQuery)
		if !ok {
			return
		}
		uq.Where(seal.In(key, val...))
	}
}

// SealULike 模糊查询
func SealULike(key string, val string) ClauseHandler {
	return func(q interface{}) {
		uq, ok := q.(*query.UpdateQuery)
		if !ok {
			return
		}
		uq.Where(seal.Like(key, val))
	}
}

// SealUOp 一般操作符 > < >= <= 等
func SealUOp(key string, op string, val interface{}) ClauseHandler {
	return func(q interface{}) {
		uq, ok := q.(*query.UpdateQuery)
		if !ok {
			return
		}
		uq.Where(seal.Op(key, op, val))
	}
}
