package storage

import (
	"context"
)

// DataHandler 数据处理
type DataHandler func(ctx context.Context, params interface{}) (interface{}, OptionStatus, error)

// 策略

// 场景1
// 1.读缓存，2.读库，3.缓存回写 （需要处理缓存穿透，击穿）

// 场景2
// 1.读缓存

// 场景3
// 1.读库(es...)

// 写场景
// 写场景比较好处理
// 场景4
// 1. 写库

// 场景5
// 1.写库，2.清空缓存

// 场景6
// 1. 写缓存，2. 写库   (双写)

// 场景7
// 1.写缓存， 2.写队列

// Do 执行
// 前一步操作的输出是下一步的输入
func Do(ctx context.Context, params interface{}, handlers ...DataHandler) (interface{}, error) {
	var err error
	if len(handlers) == 0 {
		return nil, err
	}
	var stat OptionStatus = OptionStatusContinue
	for _, fn := range handlers {
		params, stat, err = fn(ctx, params)
		if err != nil {
			return nil, err
		}
		if stat == OptionStatusBreak {
			return params, nil
		}
	}
	return params, nil
}
