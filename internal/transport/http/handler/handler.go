package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/thewolf27/image-previewer/internal/services"
)

type Handler struct {
	ctx      context.Context
	services *services.Services
}

func NewHandler(
	ctx context.Context,
	services *services.Services,
) *Handler {
	return &Handler{
		ctx:      ctx,
		services: services,
	}
}

func (h *Handler) Init(e *gin.Engine) {
	h.initResizeRoutes(e)
}
