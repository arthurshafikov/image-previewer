package services

import (
	"github.com/thewolf27/image-previewer/internal/core"
)

type ResizerService struct {
}

func NewResizerService() *ResizerService {
	return &ResizerService{}
}

func (rs *ResizerService) Resize(inp core.ResizeInput) error {
	return nil
}
