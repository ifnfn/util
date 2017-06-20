package config

// WechatInfo Server stores the hostname and port number
type WechatInfo struct {
	AppID           string `json:"AppID"`
	AppSecret       string `json:"AppSecret"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Token           string `json:"Token"`
	EncodedAESKey   string `json:"EncodedAESKey"`
}
