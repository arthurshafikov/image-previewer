package services

import (
	"image"
	"os"

	"github.com/arthurshafikov/image-previewer/internal/config"
	"github.com/arthurshafikov/image-previewer/internal/core"
)

type Resizer interface {
	ResizeFromURL(inp core.ResizeInput) (*os.File, error)
}

type Images interface {
	DownloadFromURLAndSaveImageToStorage(inp core.DownloadImageInput) (*core.Image, error)
	SaveResizedImageToStorage(imageName string, resizedImage image.Image) (*os.File, error)
}

type ImageCache interface {
	Remember(key string, callback func() (*core.Image, error)) (*core.Image, error)
	GetCachedImagesFolder() string
}

type Logger interface {
	Error(err error)
}

type Services struct {
	Resizer
}

type Deps struct {
	Config            *config.Config
	RawImageCache     ImageCache
	ResizedImageCache ImageCache
	Logger            Logger
}

func NewServices(deps Deps) *Services {
	imagesService := NewImagesService(deps.Logger, deps.RawImageCache, deps.ResizedImageCache)

	return &Services{
		Resizer: NewResizerService(deps.Logger, deps.RawImageCache, deps.ResizedImageCache, imagesService),
	}
}
