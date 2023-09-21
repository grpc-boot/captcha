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

func (wzc *WordsZhCn) Create() (dots map[int]captcha.CharDot, b64 string, thumb64 string, key string, err error) {
	return wzc.driver.Generate()
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
