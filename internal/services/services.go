package services

import "github.com/thewolf27/image-previewer/internal/core"

type Resizer interface {
	ResizeFromUrl(url string, inp core.ResizeInput) error
}

type Services struct {
	Resizer
}

func NewServices() *Services {
	return &Services{
		Resizer: NewResizerService(),
	}
}
