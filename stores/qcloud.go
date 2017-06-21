package stores

// import (
// 	"io"

// 	"roabay.com/config"

// 	cos "github.com/cxr29/cos-golang-sdk"
// )

// // QcloudStore 阿里云 OSS 存储
// type QcloudStore struct {
// 	client cos.COS
// }

// // NewQcloudStore 新建七牛云存储
// func NewQcloudStore() *QcloudStore {
// 	// return &QcloudStore{
// 	// 	client: cos.New(config.Qcloud.AppID, config.Qcloud.SecretID, config.Qcloud.SecretKey),
// 	// 	// client: cos.New(config.Qcloud.AppID, config.Qcloud.SecretID, config.Qcloud.SecretKey),
// 	// }
// }

// // Save 保存
// func (f QcloudStore) Save(key string, file io.Reader) error {
// 	F := f.client.Bucket(config.Qcloud.Bucket).Dir("").File(key)

// 	_, err := F.Upload(file, "")
// 	// ret, err := f.client.Upload(file, config.Qcloud.Bucket, key)
// 	// system.PrintInterface(ret)

// 	return err
// }

// // Get 读取数据
// func (f QcloudStore) Get(key string) (io.ReadCloser, error) {
// 	// return f.client.Download(config.Qcloud.Bucket, key)

// 	return nil, nil
// }

// // Remove 删除
// func (f QcloudStore) Remove(key string) error {
// 	// _, e := f.client.Delete(config.Qcloud.Bucket, key)

// 	// return e

// 	return nil
// }

// // Stat 读取数据
// func (f QcloudStore) Stat(key string) (Stat, error) {
// 	// res, e := f.client.StatFile(config.Qcloud.Bucket, key)

// 	// system.PrintInterface(res)

// 	return Stat{}, nil
// }

// // URL 获取资源的 URL
// func (f QcloudStore) URL(key string) string {
// 	// isPublic, _ := f.client.IsBucketPublic(config.Qcloud.Bucket)
// 	// URL := ""
// 	// if isPublic {
// 	// 	URL = f.client.GetAccessURL(config.Qcloud.Bucket, key)
// 	// } else {
// 	// 	URL = f.client.GetAccessURLWithToken(config.Qcloud.Bucket, key, 86400)
// 	// }

// 	// return URL

// 	return ""
// }

// func (f QcloudStore) List() []Stat {

// 	return nil
// }
