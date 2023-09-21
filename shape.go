package captcha

import (
	"image"
	"image/color"
	"math"
)

type shape struct {
	p image.Point
	r int
}

func (s *shape) ColorModel() color.Model {
	return color.AlphaModel
}

func (s *shape) Bounds() image.Rectangle {
	return image.Rect(s.p.X-s.r, s.p.Y-s.r, s.p.X+s.r, s.p.Y+s.r)
}

func (s *shape) At(x, y int) color.Color {
	var (
		xx, yy     = x - s.p.X, y - s.p.Y
		maxValue   = s.r * s.r
		limitValue = maxValue / 9
	)

	leftX := x - s.p.X + s.r
	leftY := y - s.p.Y - s.r/4
	topY := y - s.p.Y + s.r/3

	if xx < 0 && leftX*leftX+leftY*leftY < limitValue {
		return color.Alpha{}
	}

	if yy < -int(math.Round(float64(s.r)/3)) && xx*xx+topY*topY > limitValue {
		return color.Alpha{}
	}

	return color.Alpha{A: 255}
}
