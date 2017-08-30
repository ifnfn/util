package config

import (
	"fmt"
)

// ServerInfo stores the hostname and port number
type AccountInterface struct {
	DomainName   string `json:"DomainName"` // 域名
	Hostname     string `json:"Hostname"`   // Server name
	HTTPPort     int    `json:"HTTPPort"`   // HTTP port
	UserInfoPath string `json:"UserInfoPath"`
	UserIDPath   string `json:"UserIDPath"`
}

// HTTPBindAddress 返回REST服务器的 HTTP URL
func (s AccountInterface) GetUserInfoHTTPAddress() string {
	return fmt.Sprintf("http://%s:%d/users/%s", s.Hostname, s.HTTPPort, s.UserInfoPath)
}

func (s AccountInterface) GetUserIDHTTPAddress() string {
	return fmt.Sprintf("http://%s:%d/users/%s", s.Hostname, s.HTTPPort, s.UserIDPath)
}
