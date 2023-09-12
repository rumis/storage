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
	K string `json:"k"`
	V string `json:"p"`
}

func (p *Pair) Key() string {
	return "test_v2_pair_" + p.K
}
func (p *Pair) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}
func (p *Pair) Value(v string) error {
	err := json.Unmarshal([]byte(v), p)
	return err
}

type PairSlice []*Pair

func (ps PairSlice) ForEach(ite meta.Iterator) error {
	for k := range ps {
		err := ite(ps[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PairSlice) Value(v string) error {
	err := json.Unmarshal([]byte(v), ps)
	return err
}

type PairKey string

func (k PairKey) Key() string {
	return "test_v2_pair_" + string(k)
}

type PairKeySlice []PairKey

func (ks PairKeySlice) ForEach(ite meta.Iterator) error {
	for _, v := range ks {
		err := ite(v)
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
		K: "k1",
		V: "v1",
	}
	err = writer(ctx, kv1, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 写入多KV对象
	kv2 := &Pair{
		K: "k2",
		V: "v2",
	}
	kv3 := &Pair{
		K: "k3",
		V: "v3",
	}
	err = writer(ctx, PairSlice{
		kv2, kv3,
	}, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// 写入单对象
	writer1 := NewRedisKeyValueWriter(WithClient(rClient))
	obj1 := &Person{"obj1", 12}
	err = writer1(ctx, obj1, 0)
	if err != nil {
		t.Fatal(err)
	}
	// 写入多对象
	obj2 := &Person{"obj2", 2}
	obj3 := &Person{"obj3", 3}
	err = writer1(ctx, &Persons{obj2, obj3}, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 读取对象
	reader := NewRedisKeyValueReader(WithClient(rClient))

	// 读取单个对象-key类型为string
	rkv1 := Pair{}
	err = reader(ctx, kv1, &rkv1)
	if err != nil {
		t.Fatal(err)
	}
	if rkv1.V != kv1.V {
		t.Fatal(*kv1, rkv1)
	}

	// 读取多个KV对象
	ps := PairSlice{}
	err = reader(ctx, PairKeySlice{PairKey(kv1.K), PairKey(kv2.K), PairKey(kv3.K)}, &ps)
	if err != nil {
		t.Fatal(err)
	}
	if ps[0].V != kv1.V || ps[1].V != kv2.V || ps[2].V != kv3.V {
		t.Fatal(ps)
	}

	// 读取写入的对象值
	reader1 := NewRedisKeyValueReader(WithClient(rClient))
	expectArr := Persons{obj1, obj2, obj3}
	var allps Persons
	err = reader1(ctx, &expectArr, &allps)
	if err != nil {
		t.Fatal(err)
	}
	for i, p := range allps {
		if p.Age != expectArr[i].Age {
			t.Fatal(p)
		}
	}

}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) Key() string {
	return "test_v2_person_" + p.Name
}
func (p *Person) String() string {
	b, _ := json.Marshal(*p)
	return string(b)
}
func (p *Person) Value(v string) error {
	err := json.Unmarshal([]byte(v), p)
	return err
}

type Persons []*Person

// 遍历
func (p *Persons) ForEach(fn meta.Iterator) error {
	for _, v := range *p {
		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Persons) Value(v string) error {
	err := json.Unmarshal([]byte(v), p)
	return err
}
