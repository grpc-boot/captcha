package captcha

import (
	"image"
	"image/color"
)

type shape struct {
	p    image.Point
	maxR int
	minR int
}

func (s *shape) ColorModel() color.Model {
	return color.AlphaModel
}

func (s *shape) Bounds() image.Rectangle {
	return image.Rect(s.p.X-s.maxR, s.p.Y-s.maxR, s.p.X+s.maxR, s.p.Y+s.maxR)
}

func (s *shape) At(x, y int) color.Color {
	var (
		xx, yy   = x - s.p.X, y - s.p.Y
		value    = xx*xx + yy*yy
		minValue = s.minR * s.minR
		maxValue = s.maxR * s.maxR
	)

	if s.minR > 0 && value < minValue {
		return color.Alpha{}
	}

	if value < maxValue-5 {
		return color.Gray{
			Y: 1,
		}
	}

	if value < maxValue {
		return color.Alpha{A: 255}
	}

	return color.Alpha{}
}
