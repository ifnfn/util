package system

import (
	"bufio"
	"bytes"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// BarCodeGetPNG 生成PNG 二维码
func brCode(key string, width, height int) (barcode.Barcode, error) {
	code, err := qr.Encode(key, qr.L, qr.Unicode)
	// code, err := code39.Encode(base64)

	if err != nil {
		log.Fatal(err)
	}

	// return code, err
	return barcode.Scale(code, width, height)
}

// BarCodeGetPNG 生成PNG 二维码
func BarCodeGetPNG(key string, width, height int) (io.Reader, error) {
	code, err := brCode(key, width, height)
	if err == nil {
		b := bytes.Buffer{}
		bw := bufio.NewWriter(&b)

		png.Encode(bw, code)
		bw.Flush()
		return bufio.NewReader(&b), nil
	}

	return nil, err
}

// BarCodeGetJPEG 生成 JPEG 二维码
func BarCodeGetJPEG(key string, width, height int) (io.Reader, error) {
	code, err := brCode(key, width, height)
	if err == nil {
		b := bytes.Buffer{}
		bw := bufio.NewWriter(&b)

		jpeg.Encode(bw, code, nil)
		bw.Flush()
		return bufio.NewReader(&b), nil
	}

	return nil, err
}
