package stores

import (
	"io"
	"os"
	"path"
)

// LocalFileStore 本地文件存储
type LocalFileStore struct {
	Path string
}

// NewLocalFileStore 新建七牛云存储
func NewLocalFileStore(path string) *LocalFileStore {
	return &LocalFileStore{path}
}

// Save 保存
func (f LocalFileStore) Save(key string, file io.Reader) error {
	filePathName := path.Join(f.Path, key)

	// Destination
	dst, err := os.Create(filePathName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)

	return err
}

// Get 读
func (f LocalFileStore) Get(key string) (io.ReadCloser, error) {
	filePathName := path.Join(f.Path, key)

	return os.Open(filePathName)
}

// Remove 删除
func (f LocalFileStore) Remove(key string) error {

	return nil
}

// Stat ...
func (f LocalFileStore) Stat(key string) (Stat, error) {
	s := Stat{}

	filePathName := path.Join(f.Path, key)
	info, err := os.Stat(filePathName)
	if err == nil {
		s.Name = key
		s.Size = info.Size()
		s.UpdateTime = info.ModTime().Unix()
	}

	return s, err
}

// URL ...
func (f LocalFileStore) URL(key string) string {

	return ""
}

// List ...
func (f LocalFileStore) List() []Stat {

	return nil
}
