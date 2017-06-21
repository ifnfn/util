package config

import (
	"fmt"
	"path"
)

// ServerInfo stores the hostname and port number
type ServerInfo struct {
	DomainName     string `json:"DomainName"`     // 域名
	Hostname       string `json:"Hostname"`       // Server name
	UseHTTP        bool   `json:"UseHTTP"`        // Listen on HTTP
	UseHTTPS       bool   `json:"UseHTTPS"`       // Listen on HTTPS
	HTTPPort       int    `json:"HTTPPort"`       // HTTP port
	HTTPSPort      int    `json:"HTTPSPort"`      // HTTPS port
	CertFile       string `json:"CertFile"`       // HTTPS certificate
	KeyFile        string `json:"KeyFile"`        // HTTPS private key
	DatabaseKey    string `json:"DatabaseKey"`    // 数据库字段加密密码
	JWTSecretKey   string `json:"JWTSecretKey"`   // JWT 加密密码
	Font           string `json:"Font"`           // 验证码图片的字库
	UploadPath     string `json:"UploadPath"`     // 上传文件保留的路径
	Database       string `json:"Database"`       // 选择数据类型
	SuperClientUID string `json:"SuperClientUID"` // 超级管理员 UID
}

// HTTPBindAddress 返回REST服务器的 HTTP URL
func (s ServerInfo) HTTPBindAddress() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.HTTPPort)
}

// HTTPSBindAddress 返回服务器的 HTTPS URL
func (s ServerInfo) HTTPSBindAddress() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.HTTPSPort)
}

// HTTP 返回服务器的 Http 扩展
func (s ServerInfo) HTTP(args ...interface{}) string {
	base := fmt.Sprintf("%s:%d", s.DomainName, s.HTTPPort)

	param := fmt.Sprint(args...)
	return "http://" + path.Join(base, param)
}

// HTTPf 带格式化处理
func (s ServerInfo) HTTPf(format string, args ...interface{}) string {
	base := fmt.Sprintf("%s:%d", s.DomainName, s.HTTPPort)

	param := fmt.Sprintf(format, args...)
	return "http://" + path.Join(base, param)
}

// HTTPS 返回服务器的 Https 扩展
func (s ServerInfo) HTTPS(args string) string {
	base := fmt.Sprintf("%s:%d", s.DomainName, s.HTTPSPort)

	return "https://" + path.Join(base, args)
}

// HTTPEx ...
func (s ServerInfo) HTTPEx(port, args string) string {
	base := fmt.Sprintf("%s:%s", s.DomainName, port)

	return "http://" + path.Join(base, args)
}

// IsSuper 判断是否是超级用户
func (s ServerInfo) IsSuper(uid string) bool {
	return uid == s.SuperClientUID
}
