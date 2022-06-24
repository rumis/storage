package tutorial

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rumis/storage"
	"github.com/rumis/storage/srepo"
)

func TestPerson(t *testing.T) {

	mock := InitClient(t)

	mock.ExpectExec("INSERT INTO tal_test_person (id, name, age) VALUES (?,?,?)").WithArgs(1, "张三", 3).WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "张三", 3)
	mock.ExpectQuery("SELECT id,name,age FROM tal_test_person WHERE id=? LIMIT 1").WithArgs(1).WillReturnRows(rows)

	// 写入一条测试数据
	inster := srepo.NewSealMysqlInserter(srepo.WithDB(srepo.SealW()), srepo.WithName("tal_test_person"))
	t1 := Person{
		ID:   1,
		Name: "张三",
		Age:  3,
	}
	_, err := inster(context.TODO(), t1)
	if err != nil {
		t.Fatal(err)
	}

	// 读取
	p, err := storage.Do(context.TODO(), 1, NewNormalFlow())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)

	p2, err := storage.Do(context.TODO(), 1, NewNormalFlow())

	if err != nil {
		t.Fatal(err)
	}
	t.Log(p2)

	// 确保所有期望合格
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

}
