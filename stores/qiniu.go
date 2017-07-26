package stores

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"

	"roabay.com/util/config"
	"roabay.com/util/system"

	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
)

// QiniuStore 阿里云 OSS 存储
type QiniuStore struct {
	Client *kodo.Client
	Bucket string
	Domain string
}

// NewQiniuStore 新建七牛云存储
func NewQiniuStore(bucket, domain string) *QiniuStore {
	conf.ACCESS_KEY = config.Qiniu.AccessKey
	conf.SECRET_KEY = config.Qiniu.SecretKey

	// 创建一个Client
	c := kodo.New(0, nil)

	if bucket == "" {
		bucket = config.Qiniu.Bucket
	}
	if domain == "" {
		domain = config.Qiniu.Domain
	}

	return &QiniuStore{Client: c, Bucket: bucket, Domain: domain}
}

// Save 保存
func (f QiniuStore) Save(key string, file io.Reader) error {
	p := f.Client.Bucket(f.Bucket)

	if _, e := p.Stat(nil, key); e == nil {
		p.Delete(nil, key)
	}

	type PutRet struct {
		Hash string `json:"hash"`
		Key  string `json:"key"`
	}

	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: f.Bucket,
		//设置Token过期时间
		Expires: 3600,
	}

	// 生成一个上传token
	token := f.Client.MakeUptoken(policy)
	uploader := kodocli.NewUploader(0, nil)

	var ret PutRet

	return uploader.Put(nil, &ret, token, key, file, -1, nil)
}

// Get 读取数据
func (f QiniuStore) Get(key string) (io.ReadCloser, error) {
	baseURL := kodo.MakeBaseUrl(f.Domain, key)
	policy := kodo.GetPolicy{
		//设置Token过期时间
		Expires: 3600,
	}

	// 调用MakePrivateUrl方法返回url
	u := f.Client.MakePrivateUrl(baseURL, &policy)

	var r io.ReadCloser
	ret, err := system.HTTPGet(u, nil)
	if err == nil {
		r = ioutil.NopCloser(bytes.NewBuffer(ret))
	}

	return r, err
}

// Remove 删除
func (f QiniuStore) Remove(key string) error {
	// new一个Bucket管理对象
	p := f.Client.Bucket(f.Bucket)

	return p.Delete(nil, key)
}

// Stat 读取数据
func (f QiniuStore) Stat(key string) (Stat, error) {
	// new一个Bucket管理对象
	p := f.Client.Bucket(f.Bucket)

	s, e := p.Stat(nil, key)
	if e == nil {
		return Stat{
			Hash:       s.Hash,
			Size:       s.Fsize,
			UpdateTime: s.PutTime / 10000 / 1000,
		}, nil
	}

	return Stat{}, e
}

// URL 获取资源的 URL
func (f QiniuStore) URL(key string) string {
	baseURL := kodo.MakeBaseUrl(f.Domain, key)
	policy := kodo.GetPolicy{}

	// 调用MakePrivateUrl方法返回url
	return f.Client.MakePrivateUrl(baseURL, &policy)
}

// List 资源列表
func (f QiniuStore) List() []Stat {
	// new一个Bucket管理对象
	p := f.Client.Bucket(f.Bucket)

	var ret []Stat

	marker := ""
	for {
		items, _, markerOut, e := p.List(nil, "", "", marker, 0)
		for _, item := range items {
			ret = append(ret, Stat{
				Name:       item.Key,
				Hash:       item.Hash,
				MimeType:   item.MimeType,
				Size:       item.Fsize,
				UpdateTime: item.PutTime,
			})
		}
		marker = markerOut
		if e == io.EOF {
			break
		}
	}
	sort.Sort(StatArray(ret))
	return ret
}
