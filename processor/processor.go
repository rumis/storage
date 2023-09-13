package processor

import (
	"context"
)

// ReadProcessor 数据读取处理器
type ReadProcessor interface {
	// Read 数据读取
	// @params ctx 上下文
	// @params in 输入参数
	// @params out 输出结果
	// @params r 是否继续执行后续操作
	Read(ctx context.Context, in interface{}, out interface{}) bool
	// 数据写入
	// @params ctx 上下文
	// @params in 输入参数
	// @params out 输出结果
	// @params r 是否继续执行后续操作
	Write(ctx context.Context, in interface{}, out interface{}) bool
	// 加锁
	// 基于以下假设： 1.尝试等待一直到持有锁；2.假设1存在超时时间，超时后继续运行；3.基于假设2，不可作为独占锁使用，后续均为幂等操作
	Lock(context.Context) error
	// 释放锁
	// 基于以下假设： 1.多次释放不会导致错误
	UnLock(context.Context) error
}

// ReadProcessor 数据写入处理器
type WriteProcessor interface {
	Write(context.Context, interface{}) (interface{}, bool)
}

// TaskProcessor 常驻任务处理器
type TaskProcessor interface {
	Start(context.Context)
}
