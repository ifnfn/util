package system

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/BurntSushi/graphics-go/graphics"
)

// GraphicsScale 图片缩放
func GraphicsScale(src image.Image, w int, h int) (image.Image, error) {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	// err := graphics.Scale(m, src)
	err := graphics.Thumbnail(m, src)
	return m, err
}

// ImageScale 图片缩放, 当 w 或者 h 为 0时，按原图宽高比计算, 全为 0时，大小不变
func ImageScale(src io.Reader, dest io.Writer, w, h int) error {
	img, format, err := image.Decode(src)
	if err != nil {
		println(err.Error)
		return err
	}
	r := img.Bounds()
	rDx := r.Dx()
	rDy := r.Dy()

	if w == 0 && h == 0 {
		w = rDx
		h = rDy
	} else {
		if h == 0 {
			h = int(float32(rDy) / (float32(rDx) / float32(w)))
		} else if w == 0 {
			w = int(float32(rDx) / (float32(rDy) / float32(h)))
		}
	}
	if rDy < h {
		h = rDy
	}

	if rDx < w {
		w = rDx
	}

	if h > 5000 {
		w /= 2
		h /= 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := graphics.Thumbnail(dst, img); err == nil {
		if format == "jpeg" {
			return jpeg.Encode(dest, dst, nil)
		} else if format == "png" {
			return jpeg.Encode(dest, img, nil)
		}
	}

	return errors.New("format error")
}

// ImageScaleRatio 根据比较缩放
func ImageScaleRatio(src io.Reader, dest io.Writer, ratio float32) error {
	img, format, _ := image.Decode(src)
	r := img.Bounds()
	w := int(float32(r.Dx()) * ratio)
	h := int(float32(r.Dy()) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := graphics.Thumbnail(dst, img); err == nil {
		if format == "jpeg" {
			return jpeg.Encode(dest, dst, nil)
		} else if format == "png" {
			return jpeg.Encode(dest, img, nil)
		}
	}
	return errors.New("format error")
}

// ImageScaleWidth 指定宽度，保持宽高比
func ImageScaleWidth(src io.Reader, dest io.Writer, w int) error {
	img, format, _ := image.Decode(src)
	r := img.Bounds()
	h := int(float32(r.Dy()) * (float32(r.Dx()) / float32(w)))

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := graphics.Thumbnail(dst, img); err == nil {
		if format == "jpeg" {
			return jpeg.Encode(dest, dst, nil)
		} else if format == "png" {
			return jpeg.Encode(dest, img, nil)
		}
	}
	return errors.New("format error")
}

// ImageSubImage 图片裁剪
func ImageSubImage(src io.Reader, rect image.Rectangle, dest io.Writer) error {
	m, format, _ := image.Decode(src)

	rgbImg := m.(*image.YCbCr)
	subImg := rgbImg.SubImage(rect) // .(*image.YCbCr) //图片裁剪x0 y0 x1 y1

	if format == "jpeg" {
		return jpeg.Encode(dest, subImg, nil)
	} else if format == "png" {
		return png.Encode(dest, subImg)
	}

	return errors.New("format error")
}
