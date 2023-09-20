package captcha

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"testing"
)

func TestSlide_Create(t *testing.T) {
	s := NewSlide("")
	s.Create()
}

func TestSlide_CreateShape(t *testing.T) {
	s := NewSlide("")

	img, _, _ := openImg("./shan.png")
	rgbaImg := convertRGBA(img)
	num := 10

	width := 460
	height := 280
	r := 50

	for i := 1; i < num; i++ {
		thumbImg := s.CreateShape(rgbaImg, image.Point{X: r + rand.Intn(width-2*r), Y: r + rand.Intn(height-2*r)}, 50, 15)
		data, _ := toPng(thumbImg)
		_ = os.WriteFile(fmt.Sprintf("%d.png", i), data, 0766)
	}
}

func TestScanFiles(t *testing.T) {
	list := scanFiles("./", "png", false)
	t.Logf("list: %v", list)
}
