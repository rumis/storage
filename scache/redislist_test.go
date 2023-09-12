package scache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/meta"
)

type CacheListString string
type CacheListStringSlice []CacheListString

func (s *CacheListString) Key() string {
	return "tal_test_cache_list_x1"
}
func (s *CacheListString) String() string {
	return string(*s)
}
func (s *CacheListString) Value(v string) {
	*s = CacheListString(v)
}

func (s *CacheListStringSlice) ForEach(i meta.Iterator) error {
	for _, v := range *s {
		e1 := i(v)
		if e1 != nil {
			return e1
		}
	}
	return nil
}

func TestRedisList(t *testing.T) {

	ctx := context.TODO()

	// 启动内存Redis服务并创建Client
	server, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	// 普通字符串写入
	writer1 := NewRedisListWriter(
		WithClient(rClient),
		WithExecLogger(meta.ConsoleRedisExecLogFunc))
	err = writer1(ctx, CacheListString("t1"))
	if err != nil {
		t.Fatal(err)
	}
	err = writer1(ctx, CacheListStringSlice([]CacheListString{"t2", "t3"}))
	if err != nil {
		t.Fatal(err)
	}

	// 普通字符串读取
	reader1 := NewRedisListReader(WithClient(rClient))
	var x1 CacheListString
	err = reader1(ctx, &x1)
	if err != nil {
		t.Fatal(err)
	}
	if x1 != "t1" {
		t.Fatal("read error t1")
	}
	var x2 CacheListString
	err = reader1(ctx, &x2)
	if err != nil {
		t.Fatal(err)
	}
	if x2 != "t2" {
		t.Fatal("read error t3")
	}
	var x3 CacheListString
	err = reader1(ctx, &x3)
	if err != nil {
		t.Fatal(err)
	}
	if x3 != "t3" {
		t.Fatal("read error t3")
	}
	var xEmpty CacheListString
	err = reader1(ctx, &xEmpty)
	if err != nil {
		t.Fatal(err)
	}
	if xEmpty != "" {
		t.Fatal("read error empty")
	}
	// 一般对象写入

}

type Student struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Students []Student

// 遍历
func (s Students) ForEach(fn meta.Iterator) error {
	for _, v := range s {
		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}
