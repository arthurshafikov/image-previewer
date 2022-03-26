package services

import (
	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/core"
)

type Resizer interface {
	ResizeFromUrl(inp core.ResizeInput) error
}

type FileCache interface {
	Remember(key string, callback func() (*core.Image, error)) (*core.Image, error)
}

type Services struct {
	Resizer
}

type Deps struct {
	Config    *config.Config
	FileCache FileCache
}

func NewServices(deps Deps) *Services {
	return &Services{
		Resizer: NewResizerService(deps.Config.StorageConfig.StorageFolder, deps.FileCache),
	}
}
