package captcha

import (
	"github.com/wenlng/go-captcha/captcha"
)

type WordsZhCn struct {
	driver *captcha.Captcha
}

func (wzc *WordsZhCn) Driver() *captcha.Captcha {
	return wzc.driver
}

func (wzc *WordsZhCn) Create() (map[int]captcha.CharDot, string, string, string, error) {
	return wzc.driver.Generate()
}

func (wzc *WordsZhCn) Check() bool {
	return false
}
