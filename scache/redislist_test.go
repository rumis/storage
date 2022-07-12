package scache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/pkg/ujson"
)

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
	writer1 := NewRedisListWriter(WithClient(rClient), WithPrefix("test_list_v2_"))
	err = writer1(ctx, "t1")
	if err != nil {
		t.Fatal(err)
	}
	err = writer1(ctx, []string{"t2", "t3"})
	if err != nil {
		t.Fatal(err)
	}

	// 普通字符串读取
	reader1 := NewRedisListStringReader(WithClient(rClient), WithPrefix("test_list_v2_"))
	x1, err := reader1(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if x1 != "t1" {
		t.Fatal("read error t1")
	}

	x2, err := reader1(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if x2 != "t2" {
		t.Fatal("read error t3")
	}

	x3, err := reader1(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if x3 != "t3" {
		t.Fatal("read error t3")
	}

	xEmpty, err := reader1(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if xEmpty != "" {
		t.Fatal("read error empty")
	}

	// KV对写入
	kv1 := Pair{Key: "k1", Value: "v1"}
	kv2 := Pair{Key: "k2", Value: "v2"}
	kv3 := Pair{Key: "k2", Value: "v3"} // topic未变化，保持和kv2一致

	writer1(ctx, []Pair{kv1, kv2})
	writer1(ctx, kv3)

	rK1 := NewRedisListStringReader(WithClient(rClient), WithPrefix("test_list_v2_k1"))
	rK2 := NewRedisListStringReader(WithClient(rClient), WithPrefix("test_list_v2_k2"))

	v1, err := rK1(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v1 != kv1.Value {
		t.Fatal("read error kv1")
	}

	v2, err := rK2(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v2 != kv2.Value {
		t.Fatal("read error kv1")
	}

	v3, err := rK2(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if v3 != kv3.Value {
		t.Fatal("read error kv1")
	}

	// ForEach 对象写入
	s1 := Student{"s1", 1}
	s2 := Student{"s2", 2}
	s3 := Student{"s3", 3}
	objWriter := NewRedisListWriter(WithClient(rClient), WithKeyFn(func(i interface{}) (string, error) {
		return "test_list_object_student", nil
	}))

	objWriter(ctx, s1)
	objWriter(ctx, Students{s2, s3})

	objReader := NewRedisListStringReader(WithClient(rClient), WithPrefix("test_list_object_student"))

	outS := make(Students, 0)
	for {
		o, err := objReader(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if o == "" {
			break
		}
		var out Student
		err = ujson.Unmarshal([]byte(o), &out)
		if err != nil {
			t.Fatal(err)
		}
		outS = append(outS, out)
	}

	ins := Students{s1, s2, s3}
	for k, v := range ins {
		if v.Name != outS[k].Name || v.Grade != outS[k].Grade {
			t.Fatal(outS)
		}
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
