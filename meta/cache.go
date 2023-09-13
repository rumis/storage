package meta

import (
	"context"
	"time"
)

// KeyValueReader K-V类型读取
type KeyValueReader func(ctx context.Context, param interface{}, out interface{}) error

// KeyValueWriter K-V类型写入
type KeyValueWriter func(ctx context.Context, param interface{}, expire time.Duration) error

// KeyValueSetNX SetNX
type KeyValueSetNX func(ctx context.Context, param interface{}, expire time.Duration) bool

// KeyValueSetExp 设置KEY超时时间
type KeyValueSetExp func(ctx context.Context, param interface{}, expire time.Duration) error

// KeyValueDeleter K-V类型删除
type KeyValueDeleter func(ctx context.Context, param interface{}) error

// ListWriter List类型写入
type ListWriter func(ctx context.Context, param interface{}) error

// ListReader  List类型读取，每次读取一个值
type ListReader func(ctx context.Context, out interface{}) error
