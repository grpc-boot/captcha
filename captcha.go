package captcha

type Captcha interface {
	Check(dots string, dct map[int]Dot, span int) bool
	Create() (dots map[int]Dot, b64 string, thumb64 string, key string, err error)
}
