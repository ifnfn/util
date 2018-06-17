package stores

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"

	"github.com/ifnfn/util/config"
	"github.com/ifnfn/util/system"
)

// QiniuStore 阿里云 OSS 存储
type QiniuStore struct {
	Client *qbox.Mac
	Bucket string
	Domain string
	Config *storage.Config
}

// NewQiniuStore 新建七牛云存储
func NewQiniuStore(bucket, domain string) *QiniuStore {
	mac := qbox.NewMac(config.Qiniu.AccessKey, config.Qiniu.SecretKey)

	// // 创建一个Client
	// c := kodo.New(0, nil)

	if bucket == "" {
		bucket = config.Qiniu.Bucket
	}
	if domain == "" {
		domain = config.Qiniu.Domain
	}
	cfg := &storage.Config{
		// 空间对应的机房
		Zone: &storage.ZoneHuadong,
		// 是否使用https域名
		UseHTTPS: false,
		// 上传是否使用CDN上传加速
		UseCdnDomains: true,
	}

	return &QiniuStore{Client: mac, Bucket: bucket, Domain: domain, Config: cfg}
}

// BucketManager ...
func (f QiniuStore) BucketManager() *storage.BucketManager {
	return storage.NewBucketManager(f.Client, f.Config)
}

// Save 保存
func (f QiniuStore) Save(key string, file io.Reader) error {
	return f.SaveExt(key, file, nil)
}

// SaveExt 保存
func (f QiniuStore) SaveExt(key string, file io.Reader, ext *storage.PutExtra) error {
	if _, e := f.Stat(key); e == nil {
		f.Remove(key)
	}

	type PutRet struct {
		Hash string `json:"hash"`
		Key  string `json:"key"`
	}

	// 设置上传的策略
	putPolicy := storage.PutPolicy{
		Scope: f.Bucket,
		//设置Token过期时间
		Expires: 3600,
	}

	// 生成一个上传token
	upToken := putPolicy.UploadToken(f.Client)
	formUploader := storage.NewFormUploader(f.Config)

	var ret PutRet

	return formUploader.Put(context.Background(), &ret, upToken, key, file, -1, ext)
}

// Get 读取数据
func (f QiniuStore) Get(key string) (io.ReadCloser, error) {
	u := f.URL(key)

	var r io.ReadCloser
	ret, err := system.HTTPGet(u, nil)
	if err == nil {
		r = ioutil.NopCloser(bytes.NewBuffer(ret))
	}

	return r, err
}

// Remove 删除
func (f QiniuStore) Remove(key string) error {
	return f.BucketManager().Delete(f.Bucket, key)
}

// Stat 读取数据
func (f QiniuStore) Stat(key string) (Stat, error) {
	s, e := f.BucketManager().Stat(f.Bucket, key)
	if e == nil {
		return Stat{
			Name:       key,
			Hash:       s.Hash,
			Size:       s.Fsize,
			UpdateTime: s.PutTime / 10000 / 1000,
			MimeType:   s.MimeType,
		}, nil
	}

	return Stat{}, e
}

// URL 获取资源的 URL
func (f QiniuStore) URL(key string) string {
	// publicAccessURL := storage.MakePublicURL(f.domain, key)

	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	return storage.MakePrivateURL(f.Client, "http://"+f.Domain, key, deadline)
}

// List 资源列表
func (f QiniuStore) List() []Stat {
	// new一个Bucket管理对象
	bucketManager := f.BucketManager()

	limit := 1000
	prefix := ""
	delimiter := ""
	//初始列举marker为空
	marker := ""

	var ret []Stat

	for {
		entries, _, nextMarker, hashNext, err := bucketManager.ListFiles(f.Bucket, prefix, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}

		for _, entry := range entries {
			ret = append(ret, Stat{
				Name:       entry.Key,
				Hash:       entry.Hash,
				MimeType:   entry.MimeType,
				Size:       entry.Fsize,
				UpdateTime: entry.PutTime,
			})
		}

		if hashNext {
			marker = nextMarker
		} else {
			break
		}

	}

	sort.Sort(StatArray(ret))
	return ret
}
