package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttServerInfo stores the hostname and port number
type MqttServerInfo struct {
	Hostname         string `json:"Hostname"`
	Port             int    `json:"Port"`
	Network          string `json:"Network"`
	PingRespTimeout  int    `json:"PingRespTimeout"`
	KeepAlive        uint16 `json:"KeepAlive"`
	SecretKey        string `json:"SecretKey"`
	AdminAccessToken string `json:"AdminAccessToken"`
}

// MqttClient 客户端
type MqttClient struct {
	client   MQTT.Client
	username string
	password string
}

// mqttPassword 对密码进行加密
func mqttPassword(secretKey, password string) string {
	sum := md5.Sum([]byte(secretKey + password))

	return hex.EncodeToString(sum[:])
}

// EmqttClient 根据用户名，新一个客户端
func EmqttClient(username string) MqttClient {
	mc := MqttClient{
		username: username,
		password: mqttPassword(Mqtt.SecretKey, username),
	}

	mc.connect()

	return mc
}

func (m *MqttClient) connect() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", Mqtt.Hostname, Mqtt.Port))
	opts.SetClientID(m.username)
	opts.SetUsername(m.username)
	opts.SetPassword(m.password)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(time.Duration(Mqtt.KeepAlive) * time.Second)
	opts.SetPingTimeout(time.Duration(Mqtt.PingRespTimeout) * time.Second)

	m.client = MQTT.NewClient(opts)

	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// Disconnect 客户端断开连接
func (m *MqttClient) Disconnect() {
	m.client.Disconnect(0)
}

// Quit 客户端退出
func (m *MqttClient) Quit() {
	m.client.Disconnect(0)
}

// Publish 主题发布消息
func (m *MqttClient) Publish(topicName, message string, qos byte, retain bool) error {
	token := m.client.Publish(topicName, qos, false, message)
	token.Wait()

	return token.Error()
}

// Subscribe 主题订阅
func (m *MqttClient) Subscribe(topicFilter string, handler MQTT.MessageHandler, qos byte) error {
	token := m.client.Subscribe(topicFilter, qos, handler)
	token.Wait()
	return token.Error()
}

// Unsubscribe 主题取消订阅
func (m *MqttClient) Unsubscribe(topicFilter string) error {
	token := m.client.Unsubscribe(topicFilter)
	token.Wait()
	return token.Error()
}

// DefaultSubscribeMessage ?
func DefaultSubscribeMessage(client MQTT.Client, message MQTT.Message) {
	os.Stdout.WriteString("\n[Topic Name]\n" + string(message.Topic()) + "\n[Application Message]\n" + string(message.Payload()) + "\n")
}
