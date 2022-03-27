package services

import (
	"os"

	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/core"
)

type Resizer interface {
	ResizeFromUrl(inp core.ResizeInput) (*os.File, error)
}

type ImageCache interface {
	Remember(key string, callback func() (*core.Image, error)) (*core.Image, error)
}

type Services struct {
	Resizer
}

type Deps struct {
	Config            *config.Config
	RawImageCache     ImageCache
	ResizedImageCache ImageCache
}

func NewServices(deps Deps) *Services {
	return &Services{
		Resizer: NewResizerService(deps.Config.StorageConfig.StorageFolder, deps.RawImageCache, deps.ResizedImageCache),
	}
}
