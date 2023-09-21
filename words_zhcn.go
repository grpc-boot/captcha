package captcha

import (
	"strconv"
	"strings"

	"github.com/wenlng/go-captcha/captcha"
)

type WordsZhCn struct {
	driver *captcha.Captcha
}

func NewWordsZhCn(driver *captcha.Captcha) *WordsZhCn {
	return &WordsZhCn{
		driver: driver,
	}
}

func (wzc *WordsZhCn) Driver() *captcha.Captcha {
	return wzc.driver
}

func (wzc *WordsZhCn) Create() (dots map[int]Dot, b64 string, thumb64 string, key string, err error) {
	var charDots map[int]captcha.CharDot

	charDots, b64, thumb64, key, err = wzc.driver.Generate()
	if err != nil {
		return
	}

	dots = make(map[int]Dot, len(charDots))

	for index, dot := range charDots {
		dots[index] = Dot{
			Index:  dot.Index,
			Dx:     dot.Dx,
			Dy:     dot.Dy,
			Width:  dot.Width,
			Height: dot.Height,
		}
	}

	return
}

func (wzc *WordsZhCn) Check(dots string, dct map[int]Dot, span int) bool {
	src := strings.Split(dots, ",")
	chkRet := false

	if (len(dct) * 2) == len(src) {
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(src[j], 64)
			sy, _ := strconv.ParseFloat(src[k], 64)

			chkRet = captcha.CheckPointDistWithPadding(
				int64(sx),
				int64(sy),
				int64(dot.Dx),
				int64(dot.Dy),
				int64(dot.Width),
				int64(dot.Height),
				int64(span),
			)

			if !chkRet {
				return chkRet
			}
		}
	}

	return chkRet
}
