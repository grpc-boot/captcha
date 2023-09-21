package captcha

import "github.com/wenlng/go-captcha/captcha"

var (
	defaultOptions = func() options {
		return options{
			quality:   70,
			imgSize:   captcha.Size{Width: 300, Height: 240},
			thumbSize: captcha.Size{Width: 150, Height: 40},
		}
	}
)

type options struct {
	tp          uint8
	quality     int
	backDir     string
	maskDir     string
	driver      *captcha.Captcha
	imgSize     captcha.Size
	thumbSize   captcha.Size
	chars       []string
	backImgList []string
}

type Option func(opt *options)

func WithQuality(quality int) Option {
	return func(opt *options) {
		opt.quality = quality
	}
}

func WithDriver(driver *captcha.Captcha) Option {
	return func(opt *options) {
		opt.driver = driver
	}
}

func WithSlideBackDir(backDir string) Option {
	return func(opt *options) {
		opt.backDir = backDir
	}
}

func WithSlideMaskDir(maskDir string) Option {
	return func(opt *options) {
		opt.maskDir = maskDir
	}
}

func WithWordsZhCnThumbSize(size captcha.Size) Option {
	return func(opt *options) {
		opt.thumbSize = size
	}
}

func WithWordsZhCnImgSize(size captcha.Size) Option {
	return func(opt *options) {
		opt.imgSize = size
	}
}

func WithWordsZhCnChars(chars []string) Option {
	return func(opt *options) {
		opt.chars = chars
	}
}

func WithWordsZhBackImgList(backImgList []string) Option {
	return func(opt *options) {
		opt.backImgList = backImgList
	}
}
