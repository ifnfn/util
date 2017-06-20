package config

// AliyunInfo Server stores the hostname and port number
type AliyunInfo struct {
	Endpoint        string `json:"Endpoint"`
	AccessKeyID     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Bucket          string `json:"Bucket"`
}
