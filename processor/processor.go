package processor

import (
	"context"
	"time"
)

// ReadProcessor 数据读取处理器
type ReadProcessor interface {
	// Read 数据读取
	// @params ctx 上下文
	// @params in 输入参数
	// @params out 输出结果
	// @params r 是否继续执行后续操作
	Read(ctx context.Context, in interface{}, out interface{}, exp time.Duration) error
	// 数据写入
	// @params ctx 上下文
	// @params in 输入参数
	// @params out 输出结果
	// @params r 是否继续执行后续操作
	Write(ctx context.Context, in interface{}, exp time.Duration) error
}

// ReadProcessor 数据写入处理器
type WriteProcessor interface {
	Write(context.Context, interface{}) (interface{}, bool)
}

// TaskProcessor 常驻任务处理器
type TaskProcessor interface {
	Start(context.Context)
}
