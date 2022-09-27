package locker

import (
	"context"
	"time"
)

// DefaultExpire 默认超时时间 200ms
var DefaultExpire time.Duration = time.Millisecond * 200

// DefaultRetryTimes 默认重试次数 3次
var DefaultRetryTimes int = 3

// DefaultRetrySpan 默认重试间隔 70ms
var DefaultRetrySpan time.Duration = time.Microsecond * 70

// LockerAdder 加锁
type LockerAdder func(ctx context.Context, key string) bool

// LockerDeleter 删除锁
type LockerDeleter func(ctx context.Context, key string) error

// LockerOptionHandler 读取锁配置选项
type LockerOptionHandler func(*Locker)

// Locker 数据库读锁
type Locker struct {
	Adder      LockerAdder
	Deleter    LockerDeleter
	Expire     time.Duration
	RetryTimes int
	RetrySpan  time.Duration
}

// DefaultLocker 创建默认Locker对象
func DefaultLocker() Locker {
	return Locker{
		Adder:      func(ctx context.Context, key string) bool { return false },
		Deleter:    func(ctx context.Context, key string) error { return nil },
		Expire:     DefaultExpire,
		RetryTimes: DefaultRetryTimes,
		RetrySpan:  DefaultRetrySpan,
	}
}

// NewLocker 创建新Locker对象
func NewLocker(opts ...LockerOptionHandler) Locker {
	l := DefaultLocker()
	for _, fn := range opts {
		fn(&l)
	}
	return l
}

// WithLockerAdder 加锁
func WithLockerAdder(a LockerAdder) LockerOptionHandler {
	return func(opts *Locker) {
		opts.Adder = a
	}
}

// WithLockerDeleter 设置locker删除
func WithLockerDeleter(d LockerDeleter) LockerOptionHandler {
	return func(opts *Locker) {
		opts.Deleter = d
	}
}

// WithLockerExpire 设置locker过期时间
func WithLockerExpire(e time.Duration) LockerOptionHandler {
	return func(opts *Locker) {
		opts.Expire = e
	}
}

// WithLockerRetryTimes 设置锁重入尝试次数
func WithLockerRetryTimes(rt int) LockerOptionHandler {
	return func(opts *Locker) {
		opts.RetryTimes = rt
	}
}

// WithLockerRetrySpan 设置锁重入尝试间隔
func WithLockerRetrySpan(rs time.Duration) LockerOptionHandler {
	return func(opts *Locker) {
		opts.RetrySpan = rs
	}
}
