package config

import (
	"fmt"

	"time"

	"gopkg.in/redis.v5"
)

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

var (
	redisInfo   *RedisInfo // Database info
	RedisClient *redis.Client
)

func RedisInit(r *RedisInfo) *redis.Client {
	redisInfo = r
	fmt.Println("Redis Address:", redisInfo.Address())
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisInfo.Address(),
		Password:     redisInfo.Password, // no password set
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		DB:           0, // use default DB
	})
	//	RedisClient.FlushDb()
	return RedisClient
}
