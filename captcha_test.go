package captcha

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"testing"
)

func TestSlide_Create(t *testing.T) {
	key, x, y, thumb64, b64 := DefaultSlide.Create()
	t.Logf("key: %s x:%d y: %d", key, x, y)

	format := `<img src="%s" /> <img style="border-radius:10px;" src="%s"/>`
	_ = os.WriteFile("res.html", []byte(fmt.Sprintf(format, thumb64, b64)), 0766)
}

func BenchmarkSlide_CreateParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _, _, _ = DefaultSlide.Create()
		}
	})
}

func TestSlide_CreateShape(t *testing.T) {
	fileList := scanFiles("./back", "jpg", "png")
	width := 460
	height := 280
	r := 50
	num := 10

	for index, file := range fileList {
		img, _, _ := openImg(file)
		rgbaImg := convertRGBA(img)

		for i := 1; i < num; i++ {

			thumbImg := DefaultSlide.CreateShape(
				rgbaImg,
				image.Point{
					X: r + rand.Intn(width-2*r),
					Y: r + rand.Intn(height-2*r),
				},
				r,
			)

			data, _ := toPng(thumbImg)

			_ = os.WriteFile(fmt.Sprintf("./mask/%d_%d.png", index+1, i), data, 0766)
		}
	}
}
