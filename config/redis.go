package config

import (
	"fmt"
	"time"

	"gopkg.in/redis.v5"
)

// RedisClient ...
type RedisClient *redis.Client

// RedisInfo redis 服务器配置
type RedisInfo struct {
	Hostname string `json:"Hostname"` // Server name
	Port     int    `json:"Port"`     // HTTP port
	Password string `json:"Password"`
}

// Address Redis 服务器的地址 URL
func (s RedisInfo) Address() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

// NewRedisClient 初始化 Redis 服务器配置
func NewRedisClient() RedisClient {
	return redis.NewClient(&redis.Options{
		Addr:         Cfg.Redis.Address(),
		Password:     Cfg.Redis.Password, // no password set
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
		PoolSize:     10,
		DB:           0, // use default DB
	})
}
