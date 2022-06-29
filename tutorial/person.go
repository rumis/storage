package tutorial

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rumis/seal"
	"github.com/rumis/seal/query"
	"github.com/rumis/storage"
	"github.com/rumis/storage/meta"
	"github.com/rumis/storage/scache"
	"github.com/rumis/storage/srepo"
)

type Person struct {
	ID   int    `json:"id" seal:"id"`
	Name string `json:"name" seal:"name"`
	Age  int    `json:"age" seal:"age"`
}

func (p *Person) Zero() bool {
	return p.ID == 0
}

// 表名，缓存前缀

// NewPersonCacheReader
func NewPersonCacheReader() func(ctx context.Context, id string) (Person, error) {
	r := scache.NewRedisKeyValueObjectReader(scache.WithClient(scache.DefaultClient()), scache.WithKeyFn(func(params interface{}) (string, error) {
		p, ok := params.(Person)
		if !ok {
			return "", errors.New("not person")
		}
		return "tal_test_person_" + strconv.FormatInt(int64(p.ID), 10), nil
	}))
	return func(ctx context.Context, id string) (Person, error) {
		idi, _ := strconv.ParseInt(id, 10, 64)
		p := Person{
			ID: int(idi),
		}
		err := r(ctx, p, &p)
		if err == redis.Nil {
			return p, nil
		}
		if err != nil {
			return p, err
		}
		return p, nil
	}
}
func NewPersonCacheWriter() func(ctx context.Context, p Person, expire time.Duration) error {
	w := scache.NewRedisKeyValueWriter(scache.WithClient(scache.DefaultClient()), scache.WithKeyFn(func(params interface{}) (string, error) {
		p, ok := params.(Person)
		if !ok {
			return "", errors.New("not person")
		}
		return "tal_test_person_" + strconv.FormatInt(int64(p.ID), 10), nil
	}))
	return func(ctx context.Context, p Person, expire time.Duration) error {
		err := w(ctx, p, expire)
		return err
	}
}

func NewPersonRepoReader() func(ctx context.Context, id string) (Person, error) {
	r := srepo.NewSealMysqlOneReader(srepo.WithName("tal_test_person"), srepo.WithDB(srepo.SealR()), srepo.WithColumns([]string{"id", "name", "age"}))
	return func(ctx context.Context, id string) (Person, error) {
		var p Person
		err := r(ctx, &p, func(q interface{}) {
			sq, ok := q.(*query.SelectQuery)
			if !ok {
				return
			}
			sq.Where(seal.Eq("id", 1))
		})
		return p, err
	}
}

// 常规-缓存-数据库数据获取流程
// 读缓存-读锁-写锁-删锁-读库-写缓存
func NewNormalFlow() storage.DataHandler {
	return func(ctx context.Context, params interface{}) (interface{}, meta.OptionStatus, error) {
		id, ok := params.(int)
		if !ok {
			return nil, meta.OptionStatusBreak, errors.New("params error")
		}
		// 读取缓存
		p, err := NewPersonCacheReader()(context.TODO(), strconv.FormatInt(int64(id), 10))
		if err == nil && p.Name != "" {
			return p, meta.OptionStatusBreak, nil
		}
		// 锁
		locker := scache.DefaultRedisLocker(scache.DefaultClient(), "person")
		l, err := locker.Reader(ctx, strconv.Itoa(p.ID))
		if err == nil && l != "" {
			// 未抢到锁
			for i := 0; i < locker.RetryTimes; i++ {
				time.Sleep(locker.RetrySpan)
				p, err := NewPersonCacheReader()(context.TODO(), strconv.FormatInt(int64(id), 10))
				if err == nil && p.Name != "" {
					return p, meta.OptionStatusBreak, nil
				}
			}
		}
		// 更新锁，读库
		err = locker.Writer(ctx, strconv.Itoa(p.ID))
		if err != nil {
			fmt.Println(err)
			// return nil, 0, err
		}

		// 读库
		p1, err := NewPersonRepoReader()(context.TODO(), strconv.FormatInt(int64(id), 10))
		if err != nil {
			return nil, meta.OptionStatusBreak, err
		}

		expire := 10 * time.Second
		if p1.Name == "" {
			// 写入个空数据
			expire = 2 * time.Second
		}
		// 写缓存
		err = NewPersonCacheWriter()(context.TODO(), p1, expire)
		if err != nil {
			fmt.Println("redis write erro")
		}
		// 删除锁
		locker.Deleter(ctx, strconv.Itoa(p.ID))

		// 返回结果
		return p, meta.OptionStatusContinue, nil
	}
}
