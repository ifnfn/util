package config

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
