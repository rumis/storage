package srepo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rumis/seal"
	"github.com/rumis/seal/builder"
	"github.com/rumis/seal/query"
)

func TestRepo(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	sealDb, err := seal.OpenWithDB(db, builder.NewMysqlBuilder())
	if err != nil {
		t.Fatal(err)
	}

	// mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test_t1 (c1, c2) VALUES (?,?)").WithArgs(1, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"c1", "c2"}).AddRow(1, 3)
	mock.ExpectQuery("SELECT c1,c2 FROM test_t1 WHERE c1=? LIMIT 1").WithArgs(1).WillReturnRows(rows)

	mock.ExpectBegin()
	// update
	mock.ExpectExec("UPDATE test_t1 SET c2=? WHERE c1=?").WithArgs(5, 1).WillReturnResult(sqlmock.NewResult(0, 1))
	rows1 := sqlmock.NewRows([]string{"c1", "c2"}).AddRow(1, 5)
	mock.ExpectQuery("SELECT c1,c2 FROM test_t1 WHERE c1=? LIMIT 1").WithArgs(1).WillReturnRows(rows1)
	mock.ExpectCommit()

	// 插入数据
	ctx := context.Background()
	inster := NewSealMysqlInserter(WithDB(sealDb), WithName("test_t1"))
	t1 := T1{
		C1: 1,
		C2: 3,
	}
	_, err = inster(ctx, t1)

	if err != nil {
		t.Fatal(err)
	}
	var t2 T1
	// 读取数据
	reader := NewSealMysqlOneReader(WithDB(sealDb), WithName("test_t1"), WithColumns([]string{"c1", "c2"}))
	err = reader(ctx, &t2, func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			t.Fatal("select query not match")
			return
		}
		sq.Where(seal.Eq("c1", 1))
	})
	if err != nil {
		t.Fatal(err)
	}

	if t2.C2 != t1.C2 {
		t.Fatal(t2)
	}

	sealTx, err := sealDb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	// 更新数据
	t3 := T1{
		C2: 5,
	}
	updater := NewSealMysqlUpdater(WithTX(sealTx), WithName("test_t1"))
	_, err = updater(ctx, t3, func(q interface{}) {
		sq, ok := q.(*query.UpdateQuery)
		if !ok {
			t.Fatal("select query not match")
			return
		}
		sq.Where(seal.Eq("c1", 1))
	})
	if err != nil {
		t.Fatal(err)
	}

	// 校验更新
	reader2 := NewSealMysqlOneReader(WithTX(sealTx), WithName("test_t1"), WithColumns([]string{"c1", "c2"}))
	err = reader2(ctx, &t2, func(q interface{}) {
		sq, ok := q.(*query.SelectQuery)
		if !ok {
			t.Fatal("select query not match")
			return
		}
		sq.Where(seal.Eq("c1", 1))
	})
	if err != nil {
		t.Fatal(err)
	}
	sealTx.Commit() // 提交事务

	if t2.C2 != 5 {
		t.Fatal("update error")
	}

	// tx, _ := db.Begin()
	// _, err = tx.Exec("INSERT INTO test_t1(c1,c2) values(?,?)", 1, 3)
	// tx.Commit()

	// 确保所有期望合格
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

type T1 struct {
	C1 int `seal:"c1,omitempty"`
	C2 int `seal:"c2,omitempty"`
}
