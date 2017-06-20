package stores

import (
	"io"
)

// Stat ...
type Stat struct {
	Name       string
	Hash       string
	MimeType   string
	Size       int64
	UpdateTime int64
}

type StatArray []Stat

// Store 存储器接口
type Store interface {
	Save(key string, file io.Reader) error //
	Get(key string) (io.ReadCloser, error) //
	Remove(key string) error               //
	Stat(key string) (Stat, error)         //
	URL(key string) string
	List() []Stat
}

// Copy 将一个 Store 的数据拷贝到另一个 Store
func Copy(key string, src Store, dest Store) error {
	var err error
	var r io.ReadCloser

	if r, err = src.Get(key); err == nil {
		err = dest.Save(key, r)
	}

	return err
}

func (a StatArray) Len() int {
	return len(a)
}

func (a StatArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a StatArray) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].UpdateTime < a[i].UpdateTime
}
