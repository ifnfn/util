package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Config ...
var (
	Server   ServerInfo
	Redis    RedisInfo
	Aliyun   AliyunInfo
	InfluxDB InfluxdbInfo
	MySQL    MySQLInfo
	MongoDB  MongodbInfo
	Qiniu    QiniuInfo
	Mqtt     MqttServerInfo
	Wechat   WechatInfo
	Email    EmailInfo
	Account  AccountInterface
	// Qcloud   QcloudInfo
)

// Config contains the application settings
type Config struct {
	Server   ServerInfo       `json:"Server"`
	Redis    RedisInfo        `json:"Redis"`
	Aliyun   AliyunInfo       `json:"Aliyun"`
	InfluxDB InfluxdbInfo     `json:"InfluxDB"`
	MySQL    MySQLInfo        `json:"MySQL"`
	MongoDB  MongodbInfo      `json:"MongoDB"`
	Qiniu    QiniuInfo        `json:"Qiniu"`
	Mqtt     MqttServerInfo   `json:"Mqtt"`
	Wechat   WechatInfo       `json:"Wechat"`
	Email    EmailInfo        `json:"Email"`
	Account  AccountInterface `json:"Account"`
	// Qcloud   QcloudInfo     `json:"Qcloud"`
}

// Load 从文中加载配置
func (c *Config) Load(configFile string) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		panic(err.Error())
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
	if err != nil {
		panic(err.Error())
	}

	// Parse the config
	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		panic(fmt.Sprintf("Could not parse %q: %s", configFile, err.Error()))
	}
}

// NewConfig ...
func NewConfig(configFile string) *Config {
	config := &Config{}

	if configFile == "" {
		configFile = "config" + string(os.PathSeparator) + "config.json"
	}
	config.Load(configFile)

	Redis = config.Redis
	Aliyun = config.Aliyun
	Server = config.Server
	InfluxDB = config.InfluxDB
	MySQL = config.MySQL
	Mqtt = config.Mqtt
	MongoDB = config.MongoDB
	Qiniu = config.Qiniu
	Wechat = config.Wechat
	Email = config.Email
	Account = config.Account
	// Qcloud = config.Qcloud

	return config
}
