package handler

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(
	ctx context.Context,
) *Handler {
	return &Handler{
		ctx: ctx,
	}
}

func (h *Handler) Init(e *gin.Engine) {
	h.initResizeRoutes(e)
}
