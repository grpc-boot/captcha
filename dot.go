package captcha

type Dot struct {
	Index  int `json:"index,omitempty"`
	Dx     int `json:"dx,omitempty"`
	Dy     int `json:"dy,omitempty"`
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}
