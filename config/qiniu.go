package config

// QiniuInfo Server stores the hostname and port number
type QiniuInfo struct {
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
	Domain    string `json:"Domain"` // = "xxxx.com2.z0.glb.qiniucdn.com"
	Bucket    string `json:"Bucket"`
}
