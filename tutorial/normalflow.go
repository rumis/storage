package tutorial

import (
	"context"

	"github.com/rumis/storage"
	"github.com/rumis/storage/meta"
)

// 常规-缓存-数据库数据获取流程
func NewNormalFlow() storage.DataHandler {
	return func(ctx context.Context, params interface{}) (interface{}, meta.OptionStatus, error) {

		return nil, meta.OptionStatusContinue, nil
	}
}
