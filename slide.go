package captcha

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

var (
	// DefaultSlide default slide driver object
	DefaultSlide = NewSlide("", "", 70)
)

var (
	//go:embed mask
	maskFileDir embed.FS

	//go:embed back
	backFileDir embed.FS
)

// Slide _
type Slide struct {
	quality  int
	maskList []*image.RGBA
	backList []*image.RGBA
}

// SlideWithBack new slide with background image list
func SlideWithBack(backDir string) *Slide {
	return NewSlide(backDir, "", 70)
}

// SlideWithMask new slide with mask image list
func SlideWithMask(maskDir string) *Slide {
	return NewSlide("", maskDir, 70)
}

// NewSlide new slide with background and mask image list
func NewSlide(backDir, maskDir string, quality int) *Slide {
	s := &Slide{
		quality: quality,
	}
	s.backList = s.loadResource(backDir, &backFileDir, "back")
	s.maskList = s.loadResource(maskDir, &maskFileDir, "mask")
	return s
}

// loadResource load resource.
//
// First load from the local directory, if the loading fails, the default data will be used.
func (s *Slide) loadResource(localDir string, fileDir *embed.FS, dirName string) (rgbaList []*image.RGBA) {
	var (
		files     = make([]string, 0)
		fileTypes = []string{"png", "jpg", "jpeg"}
	)

	if localDir != "" {
		files = scanFiles(localDir, fileTypes...)
	}

	if len(files) > 0 {
		rgbaList = make([]*image.RGBA, 0, len(files))
		for _, fileName := range files {
			sp, _, _ := openImg(fileName)
			rgbaList = append(rgbaList, convertRGBA(sp))
		}

		return
	}

	fsList, _ := fileDir.ReadDir(dirName)
	rgbaList = make([]*image.RGBA, 0, len(fsList))

	for _, file := range fsList {
		for _, ft := range fileTypes {
			if strings.HasSuffix(file.Name(), ft) {
				f, _ := fileDir.Open(fmt.Sprintf("%s/%s", dirName, file.Name()))
				rgbaImg, _ := file2Rgba(f)
				rgbaList = append(rgbaList, rgbaImg)
				_ = f.Close()
				break
			}
		}
	}

	return
}

// CreateShape create shape mask.
func (s *Slide) CreateShape(img *image.RGBA, point image.Point, r int) *image.RGBA {
	sp := &shape{
		p: point,
		r: r,
	}

	thumbImg := img.SubImage(image.Rect(sp.p.X-sp.r, sp.p.Y-sp.r, sp.p.X+sp.r, sp.p.Y+sp.r))
	shapeImg := image.NewRGBA(thumbImg.Bounds())
	draw.DrawMask(shapeImg, shapeImg.Bounds(), thumbImg, thumbImg.Bounds().Min, sp, thumbImg.Bounds().Min, draw.Over)

	return shapeImg
}

// cutShape cut shape mask, return key, mask and big picture.
func (s *Slide) cutShape(img *image.RGBA, point image.Point, r int, quality int) (key string, thumb64, b64 []byte) {
	shapeImg := s.CreateShape(img, point, r)
	thumb64, _ = toBase64(MimePng, shapeImg, 0)
	key = hash(thumb64)

	rgbaImg := s.maskList[rand.Intn(len(s.maskList))]
	draw.DrawMask(img, shapeImg.Bounds(), rgbaImg, rgbaImg.Bounds().Min, rgbaImg, rgbaImg.Bounds().Min, draw.Over)

	b64, _ = toBase64(MimeJpeg, img, quality)
	return
}

// Create generate slide key, x, y, and pictures.
func (s *Slide) Create() (dots map[int]Dot, b64 string, thumb64 string, key string, err error) {
	return s.CreateCustom(50, s.quality)
}

// check 校验是否通过
func (s *Slide) check(localX, localY, reqX, reqY int, span int) bool {
	if localY != reqY {
		return false
	}

	return math.Abs(float64(localX-reqX)) <= float64(span)
}

func (s *Slide) Check(dots string, dct map[int]Dot, span int) bool {
	if len(dct) != 1 {
		return false
	}

	point, exists := dct[0]
	if !exists {
		return false
	}

	src := strings.Split(dots, ",")
	if len(src) != 2 {
		return false
	}

	reqX, _ := strconv.Atoi(src[0])
	reqY, _ := strconv.Atoi(src[1])

	return s.check(point.Dx, point.Dy, reqX, reqY, span)
}

// CreateCustom generate slide key, x, y, and pictures by custom.
//
// This method is only used for custom background and masks images
func (s *Slide) CreateCustom(r, quality int) (dots map[int]Dot, b64 string, thumb64 string, key string, err error) {
	var (
		index    = rand.Intn(len(s.backList))
		indexImg = s.backList[index]
	)

	rgbaImg := image.NewRGBA(indexImg.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), indexImg, indexImg.Bounds().Min, draw.Over)

	dots = map[int]Dot{
		0: {
			Dx: r + rand.Intn(rgbaImg.Bounds().Dx()-2*r),
			Dy: r + rand.Intn(rgbaImg.Bounds().Dy()-2*r),
		},
	}

	key, thumb64Bytes, b64Bytes := s.cutShape(rgbaImg, image.Point{X: dots[0].Dx, Y: dots[0].Dy}, r, quality)

	b64 = bytes2String(b64Bytes)
	thumb64 = bytes2String(thumb64Bytes)
	return
}
