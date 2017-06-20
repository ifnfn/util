package stores

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"time"

	"roabay.com/util/config"
	"roabay.com/util/system"

	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunStore 阿里云 OSS 存储
type AliyunStore struct {
	client *oss.Client
	bucket *oss.Bucket
}

// NewAliyunStore 新建阿里云 OSS 存储
func NewAliyunStore() *AliyunStore {
	client, err := oss.New(config.Cfg.Aliyun.Endpoint,
		config.Cfg.Aliyun.AccessKeyID,
		config.Cfg.Aliyun.AccessKeySecret)
	if err != nil {
		println(err.Error())
		return nil
	}

	bucket, e := client.Bucket(config.Cfg.Aliyun.Bucket)
	if e != nil {
		println(err.Error())
		return nil
	}

	return &AliyunStore{client, bucket}
}

// Save 保存
func (f AliyunStore) Save(key string, file io.Reader) error {
	return f.bucket.PutObject(key, file)
}

// Get 读
func (f AliyunStore) Get(key string) (io.ReadCloser, error) {
	return f.bucket.GetObject(key)
}

// Remove 删除
func (f AliyunStore) Remove(key string) error {
	return f.bucket.DeleteObject(key)
}

// Stat 获取文件状态
func (f AliyunStore) Stat(key string) (Stat, error) {
	var size int64
	var last int64
	var hash string

	meta, err := f.bucket.GetObjectMeta(key)
	if err == nil {
		system.PrintInterface(meta)
		size, _ = strconv.ParseInt(meta.Get("Content-Length"), 10, 0)

		if tx, e := time.Parse(time.RFC1123, meta.Get("Last-Modified")); e == nil {
			last = tx.Unix()
		}

		if b, e := base64.StdEncoding.DecodeString(meta.Get("Content-Md5")); e == nil {
			sum := md5.Sum(b)
			hash = hex.EncodeToString(sum[:])
		}
	}

	return Stat{
		Hash:       hash,
		Size:       size,
		UpdateTime: last,
	}, err
}

func (f AliyunStore) URL(key string) string {

	return ""
}

func (f AliyunStore) List() []Stat {

	return nil
}
