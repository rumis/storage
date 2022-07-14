package skafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// NewWriter 创建新的Kafka数据写入器
func NewWriter(opts ...KafkaWriterOptionHandler) (*kafka.Writer, Closer) {
	cfg := DefaultWriterConfig()
	for _, fn := range opts {
		fn(cfg)
	}
	w := kafka.NewWriter(*cfg)
	return w, func() error {
		return w.Close()
	}
}

// NewWriter1 创建新的Kafka写入器
func NewWriter1(opts ...KafkaWriterOptionHandler) (func(context.Context, ...kafka.Message) error, Closer) {
	w, closer := NewWriter(opts...)
	return func(ctx context.Context, msgs ...kafka.Message) error {
		err := w.WriteMessages(ctx, msgs...)
		if err != nil {
			// logger
			return err
		}
		return nil
	}, closer
}
