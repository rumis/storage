package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rumis/storage/srepo"
	"github.com/rumis/storage/test"
)

func TestOneCacheRepoReader(t *testing.T) {

	mock := test.InitClient()

	mock.ExpectExec("INSERT INTO tal_test_person (id, name, age) VALUES (?,?,?)").WithArgs(1, "张三", 3).WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "张三", 3)
	mock.ExpectQuery("SELECT id,name,age FROM tal_test_person WHERE id=? LIMIT 1").WithArgs(1).WillReturnRows(rows)

	// 写入一条测试数据
	inster := srepo.NewSealMysqlInserter(srepo.WithDB(srepo.SealW()), srepo.WithName("tal_test_person"))
	t1 := test.Person{
		ID:   1,
		Name: "张三",
		Age:  3,
	}
	_, err := inster(context.TODO(), t1)
	if err != nil {
		t.Fatal(err)
	}

	// 读取
	genericReader := NewOneCacheRepoReader("tal_test_person_", "tal_test_person", []string{"id", "name", "age"}, "person")

	var p test.Person
	err = genericReader(context.TODO(), 1, time.Second*10, &p, srepo.SealQEq("id", 1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)

	var p2 test.Person
	err = genericReader(context.TODO(), 1, time.Second*10, &p2)
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

func BenchmarkOneCacheRepoReader(b *testing.B) {

	mock := test.InitClient()

	mock.ExpectExec("INSERT INTO tal_test_person (id, name, age) VALUES (?,?,?)").WithArgs(1, "张三", 3).WillReturnResult(sqlmock.NewResult(1, 1))
	rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "张三", 3)
	mock.ExpectQuery("SELECT id,name,age FROM tal_test_person WHERE id=? LIMIT 1").WithArgs(1).WillReturnRows(rows)

	// 写入一条测试数据
	inster := srepo.NewSealMysqlInserter(srepo.WithDB(srepo.SealW()), srepo.WithName("tal_test_person"))
	t1 := test.Person{
		ID:   1,
		Name: "张三",
		Age:  3,
	}
	_, err := inster(context.TODO(), t1)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < b.N; i++ {
		// 读取
		genericReader := NewOneCacheRepoReader("tal_test_person_", "tal_test_person", []string{"id", "name", "age"}, "person")
		var p test.Person
		err = genericReader(context.TODO(), 1, time.Second*1, &p, srepo.SealQEq("id", 1))
		if err != nil {
			fmt.Println(err)
		}
	}
}
