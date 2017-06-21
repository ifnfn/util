package config

import (
	"fmt"
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
