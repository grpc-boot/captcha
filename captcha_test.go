package captcha

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestSlide_Create(t *testing.T) {
	captcha := NewWithOptions(TypeSlide)
	dots, b64, thumb64, key, _ := captcha.Create()
	t.Logf("key: %s x:%d y: %d", key, dots[0].Dx, dots[0].Dy)

	format := `<img src="%s" /> <img style="border-radius:10px;" src="%s"/>`
	_ = os.WriteFile("res.html", []byte(fmt.Sprintf(format, thumb64, b64)), 0766)

	res := captcha.Check(strconv.Itoa(dots[0].Dx)+","+strconv.Itoa(dots[0].Dy), dots, 5)

	if !res {
		t.Fatalf("want true, got %v", res)
	}
}

func TestWordsZhCn_Create(t *testing.T) {
	captcha := NewWithOptions(
		TypeWordsZhCn,
		WithWordsZhCnChars([]string{"我", "传", "打", "你", "了", "好", "魑", "魅", "魍", "魉", "方", "地", "看"}),
	)
	dots, b64, thumb64, key, _ := captcha.Create()
	t.Logf("key: %s dots:%+v", key, dots)

	format := `<img src="%s" /> <img style="border-radius:10px;" src="%s"/>`
	_ = os.WriteFile("res.html", []byte(fmt.Sprintf(format, thumb64, b64)), 0766)

	var dotsBuffer strings.Builder
	dotsBuffer.WriteString(strconv.Itoa(dots[0].Dx))
	dotsBuffer.WriteByte(',')
	dotsBuffer.WriteString(strconv.Itoa(dots[0].Dy))

	for i := 1; i < len(dots); i++ {
		dotsBuffer.WriteByte(',')
		dotsBuffer.WriteString(strconv.Itoa(dots[i].Dx))
		dotsBuffer.WriteByte(',')
		dotsBuffer.WriteString(strconv.Itoa(dots[i].Dy))
	}

	res := captcha.Check(dotsBuffer.String(), dots, 5)
	if !res {
		t.Fatalf("want true, got %v", res)
	}
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
