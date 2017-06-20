package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Cfg ...
var Cfg Config

// Config contains the application settings
type Config struct {
	Server   ServerInfo     `json:"Server"`
	Redis    RedisInfo      `json:"Redis"`
	Aliyun   AliyunInfo     `json:"Aliyun"`
	InfluxDB InfluxdbInfo   `json:"InfluxDB"`
	MySQL    MySQLInfo      `json:"MySQL"`
	MongoDB  MongodbInfo    `json:"MongoDB"`
	Qiniu    QiniuInfo      `json:"Qiniu"`
	Mqtt     MqttServerInfo `json:"Mqtt"`
	Wechat   WechatInfo     `json:"Wechat"`
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
	if configFile == "" {
		configFile = "config" + string(os.PathSeparator) + "config.json"
	}
	Cfg.Load(configFile)

	return &Cfg
}
