package captcha

import (
	"bytes"
	"encoding/base64"
	"hash/crc64"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"
)

const (
	MimePng  = `image/png`
	MimeJpg  = `image/jpg`
	MimeJpeg = `image/jpeg`
)

var (
	pngBase64Header  = []byte("data:image/png;base64,")
	jpgBase64Header  = []byte("data:image/jpg;base64,")
	jpegBase64Header = []byte("data:image/jpeg;base64,")
)

var (
	crc64Table = crc64.MakeTable(crc64.ISO)
)

func hash(data []byte) string {
	return strconv.FormatUint(crc64.Checksum(data, crc64Table), 10)
}

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

func file2Rgba(f fs.File) (rgbaImg *image.RGBA, err error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	rgbaImg = convertRGBA(img)
	return
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
	case MimeJpg:
		buf.Grow(len(dst) + len(jpgBase64Header))
		buf.Write(jpgBase64Header)
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
	case MimeJpg, MimeJpeg:
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

func scanFiles(dir string, extensionList ...string) []string {
	pathList := make([]string, 0)
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != dir {
			return filepath.SkipDir
		}

		for _, extension := range extensionList {
			if strings.HasSuffix(d.Name(), extension) {
				pathList = append(pathList, path)
				break
			}
		}

		return err
	})

	return pathList
}

func bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
