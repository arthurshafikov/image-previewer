package services

import (
	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/core"
)

type Resizer interface {
	ResizeFromUrl(inp core.ResizeInput) error
}

type Services struct {
	Resizer
}

type Deps struct {
	Config *config.Config
}

func NewServices(deps Deps) *Services {
	return &Services{
		Resizer: NewResizerService(deps.Config.StorageConfig.StorageFolder),
	}
}
