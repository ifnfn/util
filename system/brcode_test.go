package system

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrcode(t *testing.T) {
	a, e := BarCodeGetPNG("dddddddddxxxddddd", 100, 100)
	if e != nil {
		println(a)
		// println(e.Error())
	}

	d, _ := ioutil.ReadAll(a)
	ioutil.WriteFile("a.png", d, 0666)

	a, e = BarCodeGetJPEG("dddddddddxxxddddd", 100, 100)
	if e != nil {
		println(a)
		// println(e.Error())
	}

	d, _ = ioutil.ReadAll(a)
	ioutil.WriteFile("a.jpg", d, 0666)
}

func TestCrypt(t *testing.T) {
	text := GetRandomString(1000)
	src1 := bytes.NewBufferString(text)

	dst := bytes.NewBuffer(nil)
	src2 := bytes.NewBuffer(nil)
	AesEncrpyt(src1, dst, []byte("abxxxxc"), 1)
	AesDecrypt(dst, src2, []byte("abxxxxc"), 1)

	assert.Equal(t, text, string(src2.Bytes()))
}
