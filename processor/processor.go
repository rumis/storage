package processor

import "context"

// ReadProcessor 数据读取处理器
type ReadProcessor interface {
	Read(context.Context, interface{}) (interface{}, bool)
	Write(context.Context, interface{}) (interface{}, bool)
	Lock(context.Context)
}

type WriteProcessor interface {
	Write(context.Context, interface{}) (interface{}, bool)
}
