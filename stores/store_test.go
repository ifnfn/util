package stores

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"roabay.com/util/config"
)

func testStore(t *testing.T, s Store) {
	key := fmt.Sprintf("%d", time.Now().Unix())
	text := "sssssddddddddddddddddddsss"
	f := bytes.NewBufferString(text)
	if err := s.Save(key, f); err != nil {
		println(err.Error())
	}

	// if v, e := s.Get(key); e == nil {
	// 	defer v.Close()
	// 	x, _ := ioutil.ReadAll(v)
	// 	assert.Equal(t, text, string(x))
	// 	v.Close()
	// }

	// stat, _ := s.Stat(key)
	// println(stat.Size, stat.UpdateTime)

	// assert.NoError(t, s.Remove(key))
}

func TestAliyun(t *testing.T) {
	s := NewAliyunStore()
	testStore(t, s)
}
func TestQcloud(t *testing.T) {
	s := NewQcloudStore()
	testStore(t, s)
}

func TestQiniu(t *testing.T) {
	s := NewQiniuStore()
	testStore(t, s)
	s.List()
	if v, e := s.Get("1495368787"); e == nil {
		defer v.Close()
		x, _ := ioutil.ReadAll(v)
		ioutil.WriteFile("a.jpg", x, 0666)
	} else {
		println(e.Error())
	}
}

const confFile = "/Users/zhuzhg/works/golang/src/phoenix/config/config.json"

func init() {
	config.Init(confFile)
}
