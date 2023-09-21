package captcha

import "github.com/wenlng/go-captcha/captcha"

const (
	TypeSlide     = 1
	TypeWordsZhCn = 2
)

func NewWithOptions(tp uint8, opts ...Option) Captcha {
	opt := defaultOptions()
	opt.tp = tp

	for _, option := range opts {
		option(&opt)
	}

	switch opt.tp {
	case TypeWordsZhCn:
		if opt.driver == nil {
			opt.driver = captcha.NewCaptcha()
		}

		if len(opt.chars) > 0 {
			_ = opt.driver.SetRangChars(opt.chars)
		}

		if len(opt.backImgList) > 0 {
			opt.driver.SetBackground(opt.backImgList, true)
		}

		opt.driver.SetImageQuality(opt.quality)
		opt.driver.SetThumbSize(opt.thumbSize)
		opt.driver.SetImageSize(opt.imgSize)

		return NewWordsZhCn(opt.driver)
	}

	return NewSlide(opt.backDir, opt.maskDir, opt.quality)
}
