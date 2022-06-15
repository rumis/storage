package scache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"

	"github.com/rumis/storage"
	"github.com/rumis/storage/pkg/ujson"
)

func TestRedisKV(t *testing.T) {

	ctx := context.TODO()

	// 启动内存Redis服务并创建Client
	server, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	writer := NewRedisKeyValueWriter(WithClient(rClient), WithPrefix("test_v1_"))

	// 写入单KV对象
	kv1 := Pair{
		Key:   "k1",
		Value: "v1",
	}
	err = writer(ctx, kv1, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 写入多KV对象
	kv2 := Pair{
		Key:   "k2",
		Value: "v2",
	}
	kv3 := Pair{
		Key:   "k3",
		Value: "v3",
	}
	err = writer(ctx, []Pair{
		kv2, kv3,
	}, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// 写入单对象
	writer1 := NewRedisKeyValueWriter(WithClient(rClient), WithKeyFn(func(item interface{}) (string, error) {
		p, ok := item.(Person)
		if !ok {
			return "", ErrKeyGenerate
		}
		return "test_v1_" + p.Name, nil
	}))
	obj1 := Person{"obj1", 12}
	err = writer1(ctx, obj1, 0)
	if err != nil {
		t.Fatal(err)
	}
	// 写入多对象
	obj2 := Person{"obj2", 2}
	obj3 := Person{"obj3", 3}
	err = writer1(ctx, Persons{obj2, obj3}, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 读取对象
	reader := NewRedisKeyValueReader(WithClient(rClient), WithPrefix("test_v1_"))

	// 读取单个对象-key类型为string
	res, err := reader(ctx, kv1.Key)
	if err != nil {
		t.Fatal(err)
	}
	v1, ok := res.(string)
	if !ok {
		t.Fatal("redis read error:", res)
	}
	if v1 != kv1.Value {
		t.Fatal(kv1, v1)
	}

	// 读取多个KV对象
	tRes, err := reader(ctx, []string{kv1.Key, kv2.Key, kv3.Key})
	if err != nil {
		t.Fatal(err)
	}
	vals, ok := tRes.([]string)
	if !ok {
		t.Fatal("redis read error:", res)
	}
	if vals[0] != kv1.Value || vals[1] != kv2.Value || vals[2] != kv3.Value {
		t.Fatal(vals)
	}

	// 读取写入的对象值
	reader1 := NewRedisKeyValueReader(WithClient(rClient), WithKeyFn(func(item interface{}) (string, error) {
		p, ok := item.(Person)
		if !ok {
			return "", ErrKeyGenerate
		}
		return "test_v1_" + p.Name, nil
	}))
	expectArr := Persons{obj1, obj2, obj3}
	objRes, err := reader1(ctx, expectArr)
	if err != nil {
		t.Fatal(err)
	}
	objVals, ok := objRes.([]string)
	if !ok {
		t.Fatal("redis read error:", res)
	}

	for i, v := range objVals {
		var p Person
		err = ujson.Unmarshal([]byte(v), &p)
		if err != nil {
			t.Fatal(err)
		}
		if p.Age != expectArr[i].Age {
			t.Fatal(p)
		}
	}

	// 读取写入值 - 自定义读取Reader
	r1 := NewPersonReader(WithClient(rClient), WithKeyFn(func(item interface{}) (string, error) {
		p, ok := item.(Person)
		if !ok {
			return "", ErrKeyGenerate
		}
		return "test_v1_" + p.Name, nil
	}))
	personRes, err := r1(ctx, expectArr)
	if err != nil {
		t.Fatal(err)
	}
	ps, ok := personRes.([]Person)
	if !ok {
		t.Fatal("redis read error:", res)
	}
	for i, v := range ps {
		if v.Age != expectArr[i].Age {
			t.Fatal(v)
		}
	}

}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Persons []Person

// 遍历
func (p Persons) ForEach(fn storage.Iterator) error {
	for _, v := range p {
		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewPersonReader(hands ...RedisOptionHandler) RedisKeyValueReader {
	kvReader := NewRedisKeyValueReader(hands...)
	return func(ctx context.Context, params interface{}) (interface{}, error) {
		items, err := kvReader(ctx, params)
		if err != nil {
			return nil, err
		}
		switch v := items.(type) {
		case string:
			var p Person
			err := ujson.Unmarshal([]byte(v), &p)
			if err != nil {
				return nil, err
			}
			return p, nil
		case []string:
			ps := make([]Person, 0)
			for _, pv := range v {
				var p Person
				err := ujson.Unmarshal([]byte(pv), &p)
				if err != nil {
					return nil, err
				}
				ps = append(ps, p)
			}
			return ps, nil
		default:
			return nil, errors.New("reader return value error")
		}
	}
}
