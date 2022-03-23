package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thewolf27/image-previewer/internal/core"
)

func (h *Handler) initResizeRoutes(e *gin.Engine) {
	resize := e.Group("/resize")
	{
		resize.POST("", h.resize)
	}
}

func (h *Handler) resize(ctx *gin.Context) {
	h.services.Resizer.Resize(core.ResizeInput{})
	ctx.JSON(http.StatusOK, "")
}
