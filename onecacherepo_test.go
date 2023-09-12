package storage

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rumis/storage/v2/locker"
	"github.com/rumis/storage/v2/scache"
	"github.com/rumis/storage/v2/srepo"
	"github.com/rumis/storage/v2/test"
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
	// genericReader := NewOneCacheRepoReader("tal_test_person_", "tal_test_person", []string{"id", "name", "age"}, "person")

	genericReader := NewOneCacheRepoReader(NewOneCacheRepoOptions(
		WithCacheReader(NewOneCacheReader()),
		WithCacheWriter(NewOneCacheWriter("tal_test_person_")),
		WithRepoReader(NewOneRepoReader("tal_test_person", []string{"id", "name", "age"})),
		WithLocker(NewDefaultLocker("person")),
	))

	var p test.Person
	err = genericReader(context.TODO(), 1, time.Second*10, &p)
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

// // NewOneCacheReader 缓存对象读取
func NewOneCacheReader() scache.RedisKeyValueReader {
	r := scache.NewRedisKeyValueReader(scache.WithClient(scache.DefaultClient()))
	return r
}

// NewOneCacheWriter 缓存对象写入
func NewOneCacheWriter(prefix string) scache.RedisKeyValueWriter {
	return scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()))
	// return func(ctx context.Context, kv scache.Pair, expire time.Duration) error {
	// 	err := w(ctx, kv, expire)
	// 	return err
	// }
}

// NewOneRepoReader 通用数据库读取
func NewOneRepoReader(tableName string, columns []string) srepo.RepoGroupReader {
	r := srepo.NewSealMysqlOneReader(srepo.WithName(tableName), srepo.WithDB(srepo.SealR()), srepo.WithColumns(columns))
	return func(ctx context.Context, out interface{}, params interface{}) error {
		pstr := fmt.Sprint(params)
		id, err := strconv.Atoi(pstr)
		if err != nil {
			return err
		}
		err = r(ctx, out, srepo.SealQEq("id", id))
		return err
	}
}

// NewDefaultLocker 默认通用锁
func NewDefaultLocker(biz string) locker.Locker {
	return scache.DefaultRedisLocker(scache.DefaultClient(), biz)
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
		genericReader := NewOneCacheRepoReader(NewOneCacheRepoOptions(
			WithCacheReader(NewOneCacheReader()),
			WithCacheWriter(NewOneCacheWriter("tal_test_person_")),
			WithRepoReader(NewOneRepoReader("tal_test_person", []string{"id", "name", "age"})),
			WithLocker(NewDefaultLocker("person")),
		))
		var p test.Person
		err = genericReader(context.TODO(), 1, time.Second*1, &p)
		if err != nil {
			fmt.Println(err)
		}
	}
}
