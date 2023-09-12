package scache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/rumis/storage/v2/meta"
)

type Pair struct {
	k string
	v string
}

func (p *Pair) Key() string {
	return p.k
}
func (p *Pair) String() string {
	return p.v
}
func (p *Pair) Value(v string) error {
	p.v = v
	return nil
}

type PairSlice []*Pair

func (ps PairSlice) ForEach(ite meta.Iterator) error {
	for k, _ := range ps {
		err := ite(&ps[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func TestRedisKV(t *testing.T) {

	ctx := context.WithValue(context.Background(), meta.DefaultTraceKey, "asdf")

	// 启动内存Redis服务并创建Client
	server, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	rClient := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	writer := NewRedisKeyValueWriter(
		WithClient(rClient),
		WithExecLogger(meta.ConsoleRedisExecLogFunc))

	// 写入单KV对象
	kv1 := &Pair{
		k: "k1",
		v: "v1",
	}
	err = writer(ctx, kv1, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 写入多KV对象
	kv2 := &Pair{
		k: "k2",
		v: "v2",
	}
	kv3 := &Pair{
		k: "k3",
		v: "v3",
	}
	err = writer(ctx, PairSlice{
		kv2, kv3,
	}, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// 写入单对象
	writer1 := NewRedisKeyValueWriter(WithClient(rClient))
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
	reader := NewRedisKeyValueReader(WithClient(rClient))

	// 读取单个对象-key类型为string
	rkv1 := Pair{}
	err = reader(ctx, StringKey(kv1.k), &rkv1)
	if err != nil {
		t.Fatal(err)
	}
	if rkv1.v != kv1.v {
		t.Fatal(*kv1, rkv1)
	}

	// 读取多个KV对象
	ps := PairSlice{}
	err = reader(ctx, StringKeySlice{StringKey(kv1.k), StringKey(kv2.k), StringKey(kv3.k)}, &ps)
	if err != nil {
		t.Fatal(err)
	}
	if ps[0].v != kv1.v || ps[1].v != kv2.v || ps[2].v != kv3.v {
		t.Fatal(ps)
	}

	// // 读取写入的对象值
	// reader1 := NewRedisKeyValueObjectReader(WithClient(rClient), WithKeyFn(func(item interface{}) (string, error) {
	// 	p, ok := item.(Person)
	// 	if !ok {
	// 		return "", ErrKeyGenerate
	// 	}
	// 	return "test_v1_" + p.Name, nil
	// }))
	// expectArr := Persons{obj1, obj2, obj3}
	// var allps []Person
	// err = reader1(ctx, expectArr, &allps)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// for i, p := range allps {
	// 	if p.Age != expectArr[i].Age {
	// 		t.Fatal(p)
	// 	}
	// }

	// // 读取写入值 - 自定义读取Reader
	// r1 := NewPersonReader(WithClient(rClient), WithKeyFn(func(item interface{}) (string, error) {
	// 	p, ok := item.(Person)
	// 	if !ok {
	// 		return "", ErrKeyGenerate
	// 	}
	// 	return "test_v1_" + p.Name, nil
	// }))
	// personRes, err := r1(ctx, expectArr)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ps, ok := personRes.([]Person)
	// if !ok {
	// 	t.Fatal("redis read error:", res)
	// }
	// for i, v := range ps {
	// 	if v.Age != expectArr[i].Age {
	// 		t.Fatal(v)
	// 	}
	// }

}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) Key() string {
	return "test_v1_" + p.Name
}
func (p *Person) String() string {
	b, _ := json.Marshal(*p)
	return string(b)
}
func (p *Person) Value(v string) error {
	err := json.Unmarshal([]byte(v), p)
	return err
}

type Persons []Person

// 遍历
func (p Persons) ForEach(fn meta.Iterator) error {
	for _, v := range p {
		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// func NewPersonReader(hands ...RedisOptionHandler) RedisKeyValueReader {
// 	kvReader := NewRedisKeyValueReader(hands...)
// 	return func(ctx context.Context, params interface{}) (interface{}, error) {
// 		items, err := kvReader(ctx, params)
// 		if err != nil {
// 			return nil, err
// 		}
// 		switch v := items.(type) {
// 		case string:
// 			var p Person
// 			err := ujson.Unmarshal([]byte(v), &p)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return p, nil
// 		case []string:
// 			ps := make([]Person, 0)
// 			for _, pv := range v {
// 				var p Person
// 				err := ujson.Unmarshal([]byte(pv), &p)
// 				if err != nil {
// 					return nil, err
// 				}
// 				ps = append(ps, p)
// 			}
// 			return ps, nil
// 		default:
// 			return nil, errors.New("reader return value error")
// 		}
// 	}
// }
