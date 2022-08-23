package meta

import (
	"context"
	"fmt"
	"time"

	"github.com/rumis/seal/options"
)

var DefaultTraceKey = options.DefaultTraceKey

// RedisExecLogFunc is called each time when a redis command is executed.
// The "ts" parameter gives the time that the command takes to execute,
// args refer the exec params
// err refer to the result of the execution.
type RedisExecLogFunc func(ctx context.Context, ts time.Duration, args interface{}, err error)

// ConsoleRedisExecLogFunc print the message to console
func ConsoleRedisExecLogFunc(ctx context.Context, ts time.Duration, args interface{}, err error) {
	traceId := ctx.Value(DefaultTraceKey)
	fmt.Printf("redis command log: \n  trace:%v\n  args:%+v \n  timespan:%dns \n  error:%+v\n \n", traceId, args, ts.Nanoseconds(), err)
}

// KafkaLoggerFunc is called each time when read or write a message

// 	msg: 消息
// 	args: 参数
type KafkaLoggerFunc func(msg string, args ...interface{})
