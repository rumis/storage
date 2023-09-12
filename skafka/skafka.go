package skafka

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/rumis/storage/v2/meta"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var defaultDialer *kafka.Dialer
var dialers Dialers

type Dialers struct {
	m  sync.Mutex
	di map[string]*kafka.Dialer
	o  sync.Once
}

// Get 获取链接器
func (d *Dialers) Get(name string) (*kafka.Dialer, bool) {
	d.m.Lock()
	defer d.m.Unlock()
	dialer, ok := d.di[name]
	return dialer, ok
}

// Set 存储一个链接器
func (d *Dialers) Set(name string, dialer *kafka.Dialer) {
	d.o.Do(func() {
		d.di = make(map[string]*kafka.Dialer)
	})
	d.m.Lock()
	defer d.m.Unlock()
	d.di[name] = dialer
}

// KafkaDialerOptionHandler Kafka链接器配置选项
type KafkaDialerOptionHandler func(*kafka.Dialer)

// KafkaOptionHandler Kafka配置选项
type KafkaReaderOptionHandler func(*kafka.ReaderConfig)

// KafkaOptionHandler Kafka配置选项
type KafkaWriterOptionHandler func(*kafka.WriterConfig)

type Closer func() error

// D_WithUserNamePassword 设置用户名称密码
func D_WithUserNamePassword(uname string, password string) KafkaDialerOptionHandler {
	return func(d *kafka.Dialer) {
		if uname == "" {
			return
		}
		mechanism := plain.Mechanism{
			Username: uname,
			Password: password,
		}
		d.SASLMechanism = mechanism
	}
}

// D_WithCA 设置证书
func D_WithCA(ca string) KafkaDialerOptionHandler {
	return func(d *kafka.Dialer) {
		if ca == "" {
			return
		}
		caCert := x509.NewCertPool()
		caCert.AppendCertsFromPEM([]byte(ca))
		d.TLS = &tls.Config{
			RootCAs:            caCert,
			InsecureSkipVerify: true,
		}
	}
}

// InitDefaultDialer 初始化默认的Dialer
func InitDefaultDialer(opts ...KafkaDialerOptionHandler) {
	dname := "storage_kafka_default"
	NewDialer(dname, opts...)
	dialer, _ := dialers.Get(dname)
	defaultDialer = dialer
}

// NewDialer 创建新的拨号器
func NewDialer(name string, opts ...KafkaDialerOptionHandler) {
	dialer := &kafka.Dialer{}
	for _, fn := range opts {
		fn(dialer)
	}
	dialers.Set(name, dialer)
}

// GetDialer 获取链接
func GetDialer(name string) (*kafka.Dialer, bool) {
	return dialers.Get(name)
}

// DefaultDialer 获取默认的链接对象
func DefaultDialer() *kafka.Dialer {
	return defaultDialer
}

// 读取器配置设置
// DefaultReaderConfig 默认读取器配置
func DefaultReaderConfig() *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Dialer:   DefaultDialer(),
		MaxBytes: 1e7,
		MinBytes: 1e3,
	}
}

// R_WithDialer 读取器 配置链接器
func R_WithDialer(dialer *kafka.Dialer) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.Dialer = dialer
	}
}

// R_WithBroker 读取器 配置Broker
func R_WithBrokers(addrs []string) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.Brokers = addrs
	}
}

// R_WithGroup 读取器 配置消费者组
func R_WithGroupID(group string) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.GroupID = group
	}
}

// R_WithTopic 读取器 配置主题
func R_WithTopic(topic string) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.Topic = topic
	}
}

// R_WithMinBytes 读取器 配置最少读取字节数
func R_WithMinBytes(bcnt int) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.MinBytes = bcnt
	}
}

// R_WithMaxBytes 读取器 配置最多读取字节数
func R_WithMaxBytes(bcnt int) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.MaxBytes = bcnt
	}
}

// R_WithLogger 读取器 日志
func R_WithLogger(fn meta.KafkaLoggerFunc) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.Logger = kafka.LoggerFunc(fn)
	}
}

// R_WithErrorLogger 读取器 错误日志
func R_WithErrorLogger(fn meta.KafkaLoggerFunc) KafkaReaderOptionHandler {
	return func(rc *kafka.ReaderConfig) {
		rc.ErrorLogger = kafka.LoggerFunc(fn)
	}
}

// 写入器配置设置
// DefaultWriterConfig 默认写入器器配置
func DefaultWriterConfig() *kafka.WriterConfig {
	return &kafka.WriterConfig{
		Dialer:     DefaultDialer(),
		BatchBytes: 1e7,
	}
}

// W_WithDialer 写入器 配置链接器
func W_WithDialer(dialer *kafka.Dialer) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.Dialer = dialer
	}
}

// W_WithBroker 写入器 配置Broker
func W_WithBrokers(addrs []string) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.Brokers = addrs
	}
}

// W_WithTopic 写入器 配置主题
func W_WithTopic(topic string) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.Topic = topic
	}
}

// W_WithBatchBytes 写入器 配置批量发送字节数上限
func W_WithBatchBytes(bcnt int) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.BatchBytes = bcnt
	}
}

// W_WithLogger 写入器 日志
func W_WithLogger(fn meta.KafkaLoggerFunc) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.Logger = kafka.LoggerFunc(fn)
	}
}

// W_WithErrorLogger 写入器 错误日志
func W_WithErrorLogger(fn meta.KafkaLoggerFunc) KafkaWriterOptionHandler {
	return func(wc *kafka.WriterConfig) {
		wc.ErrorLogger = kafka.LoggerFunc(fn)
	}
}
