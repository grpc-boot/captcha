package captcha

type Captcha interface {
	Check(dots string, dct map[int]Dot, span int) bool
}
