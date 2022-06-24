package locker

import (
	"context"
	"time"
)

var DefaultExpire time.Duration = time.Millisecond * 200
var DefaultRetryTimes int = 3
var DefaultRetrySpan time.Duration = time.Microsecond * 70

type LockerWriter func(ctx context.Context, key string) error
type LockerReader func(ctx context.Context, key string) (string, error)

// LockerOptionHandler 数据库配置选项
type LockerOptionHandler func(*Locker)

// Locker 数据库读锁
type Locker struct {
	Reader     LockerReader
	Writer     LockerWriter
	Expire     time.Duration
	RetryTimes int
	RetrySpan  time.Duration
}

// DefaultLocker 创建默认Locker对象
func DefaultLocker() Locker {
	return Locker{
		Writer:     func(ctx context.Context, key string) error { return nil },
		Reader:     func(ctx context.Context, key string) (string, error) { return "", nil },
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

// WithLockerWriter 设置locker写入器
func WithLockerWriter(w LockerWriter) LockerOptionHandler {
	return func(opts *Locker) {
		opts.Writer = w
	}
}

// WithLockerReader 设置locker读取器
func WithLockerReader(r LockerReader) LockerOptionHandler {
	return func(opts *Locker) {
		opts.Reader = r
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
