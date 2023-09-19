package captcha

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

const (
	MimePng  = `image/png`
	MimeJpeg = `image/jpeg`
)

var (
	pngBase64Header  = []byte("data:image/png;base64,")
	jpegBase64Header = []byte("data:image/jpeg;base64,")
)

func openImg(fileName string) (img image.Image, fileType string, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	return image.Decode(f)
}

func convertRGBA(img image.Image) *image.RGBA {
	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, img.Bounds(), img, image.Point{}, draw.Over)
	return rgbaImg
}

func convertBase64(mimeType string, data []byte) []byte {
	encodeLength := base64.StdEncoding.EncodedLen(len(data))
	dst := make([]byte, encodeLength)
	base64.StdEncoding.Encode(dst, data)

	buf := bytes.NewBuffer(nil)

	switch mimeType {
	case MimePng:
		buf.Grow(len(dst) + len(pngBase64Header))
		buf.Write(pngBase64Header)
	case MimeJpeg:
		buf.Grow(len(dst) + len(jpegBase64Header))
		buf.Write(jpegBase64Header)
	default:
		return nil
	}

	buf.Write(dst)
	return buf.Bytes()
}

func toBase64(mimeType string, img image.Image, quality int) (dst []byte, err error) {
	var (
		data []byte
	)

	switch mimeType {
	case MimePng:
		data, err = toPng(img)
	case MimeJpeg:
		data, err = toJpg(img, quality)
	default:
		err = ErrFormat
		return
	}

	dst = convertBase64(mimeType, data)

	return
}

func toJpg(img image.Image, quality int) (data []byte, err error) {
	out := bytes.NewBuffer(nil)
	if quality < 1 {
		quality = jpeg.DefaultQuality
	}

	err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return
	}

	data = out.Bytes()
	return
}

func toPng(img image.Image) (data []byte, err error) {
	out := bytes.NewBuffer(nil)
	err = png.Encode(out, img)
	if err != nil {
		return
	}

	data = out.Bytes()
	return
}

func savePng(fileName string, img image.Image) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return
	}

	defer func() {
		_ = outFile.Close()
	}()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		return
	}

	return b.Flush()
}

func saveJpg(fileName string, img image.Image, quality int) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return
	}

	defer func() {
		_ = outFile.Close()
	}()

	b := bufio.NewWriter(outFile)

	if quality < 1 {
		quality = jpeg.DefaultQuality
	}

	err = jpeg.Encode(b, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return
	}

	return b.Flush()
}

func cutShap(img image.Image) {

}
