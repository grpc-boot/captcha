package captcha

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"math/rand"
	"os"
	"strings"
)

var (
	//go:embed masks
	f embed.FS
)

type Slide struct {
	maskList []*image.RGBA
}

func NewSlide(maskDir string) *Slide {
	s := &Slide{
		maskList: make([]*image.RGBA, 0, 32),
	}

	var files = make([]string, 0)
	if maskDir != "" {
		files = scanFiles(maskDir, "png", false)
	}

	if len(files) < 1 {
		fsList, _ := f.ReadDir("masks")
		for _, file := range fsList {
			if strings.HasSuffix(file.Name(), "png") {
				files = append(files, "masks/"+file.Name())
			}
		}
	}

	for _, fileName := range files {
		sp, _, _ := openImg(fileName)
		s.maskList = append(s.maskList, convertRGBA(sp))
	}

	return s
}

func (s *Slide) CreateShape(img *image.RGBA, point image.Point, maxR, minR int) *image.RGBA {
	sp := &shape{
		p:    point,
		maxR: maxR,
		minR: minR,
	}

	thumbImg := img.SubImage(image.Rect(sp.p.X-sp.maxR, sp.p.Y-sp.maxR, sp.p.X+sp.maxR, sp.p.Y+sp.maxR))
	shapeImg := image.NewRGBA(thumbImg.Bounds())
	draw.DrawMask(shapeImg, shapeImg.Bounds(), thumbImg, thumbImg.Bounds().Min, sp, thumbImg.Bounds().Min, draw.Over)

	return shapeImg
}

func (s *Slide) CutShape(img *image.RGBA, point image.Point, maxR, minR int) (thumb64, b64 []byte) {
	shapeImg := s.CreateShape(img, point, maxR, minR)
	thumb64, _ = toBase64(MimePng, shapeImg, 0)
	rgbaImg := s.maskList[rand.Intn(len(s.maskList))]
	draw.DrawMask(img, shapeImg.Bounds(), rgbaImg, rgbaImg.Bounds().Min, rgbaImg, rgbaImg.Bounds().Min, draw.Over)
	b64, _ = toBase64(MimeJpeg, img, 90)

	return
}

func (s *Slide) Create() {
	img, _, _ := openImg("./shan.png")
	rgbaImg := convertRGBA(img)

	width := 460
	height := 280
	r := 50

	x := r + rand.Intn(width-2*r)
	y := r + rand.Intn(height-2*r)

	thumb64, b64 := s.CutShape(rgbaImg, image.Point{X: x, Y: y}, r, 15)
	format := `<img src="%s" /> <img style="border-radius:10px;" src="%s"/>`

	_ = os.WriteFile("res.html", []byte(fmt.Sprintf(format, thumb64, b64)), 0766)
}
