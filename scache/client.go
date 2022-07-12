package scache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config Redis配置
type Config struct {
	Name     string // instance name
	Addr     string // host:port address.
	Username string // username
	Password string // password
	DB       int    // selected db
	PoolSize int    // connection pool size, must > 3
}

var inst *redis.Client // 默认链接
var servers sync.Map   // 链接池

// DefaultClient 获取默认的客户端
func DefaultClient() *redis.Client {
	return inst
}

// SetDefaultClient 设置默认链接
// 未加锁，需要程序启动时初始化
func SetDefaultClient(c *redis.Client) {
	inst = c
}

// GetClient 获取Redis池链接中的链接
func GetClient(name string) (*redis.Client, error) {
	s, ok := servers.Load(name)
	if !ok {
		return nil, fmt.Errorf("server name [%s] not found", name)
	}
	c, ok := s.(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("server name [%s] is a illegal *redis.Client", name)
	}
	return c, nil
}

// NewClient 初始化 创建新的Redis链接
// 默认会加入连接池
func NewClient(cfg Config) (*redis.Client, error) {
	if cfg.PoolSize < 3 {
		cfg.PoolSize = 3
	}
	// init redis client
	newc := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  2 * time.Second, // default read and write timeout is 2s
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 3,
	})
	res, err := newc.Ping(context.Background()).Result()
	if err != nil || res != "PONG" {
		return nil, fmt.Errorf("ping redis [%s] failed, error:%s", cfg.Addr, err.Error())
	}
	servers.Store(cfg.Name, newc)
	return newc, nil
}

// NewDefaultClient 创建默认的Redis链接
func NewDefaultClient(cfg Config) error {
	defaultClient, err := NewClient(cfg)
	if err != nil {
		return err
	}
	SetDefaultClient(defaultClient)
	return nil
}
