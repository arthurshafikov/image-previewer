package services

import (
	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/core"
)

type Resizer interface {
	ResizeFromUrl(inp core.ResizeInput) error
}

type ImageCache interface {
	Remember(key string, callback func() (*core.Image, error)) (*core.Image, error)
}

type Services struct {
	Resizer
}

type Deps struct {
	Config     *config.Config
	ImageCache ImageCache
}

func NewServices(deps Deps) *Services {
	return &Services{
		Resizer: NewResizerService(deps.Config.StorageConfig.StorageFolder, deps.ImageCache),
	}
}
