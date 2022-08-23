package skafka

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"
)

// NewReader 创建新的Kafka读取对象
func NewReader(opts ...KafkaReaderOptionHandler) (*kafka.Reader, Closer) {
	cfg := DefaultReaderConfig()
	for _, fn := range opts {
		fn(cfg)
	}
	r := kafka.NewReader(*cfg)
	return r, func() error {
		return r.Close()
	}
}

// NewReaderChannel 创建新的读取器 并将读取内容输出到管道
func NewReaderChannel(opts ...KafkaReaderOptionHandler) (chan kafka.Message, Closer) {
	msgCh := make(chan kafka.Message)
	r, _ := NewReader(opts...)
	ctx := context.TODO()
	go func() {
		for {
			m, err := r.ReadMessage(ctx)
			if err == io.EOF {
				continue
			}
			if err != nil {
				continue
			}
			msgCh <- m
		}
	}()
	return msgCh, func() error {
		close(msgCh)
		return r.Close()
	}
}
