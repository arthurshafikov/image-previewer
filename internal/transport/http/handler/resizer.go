package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initResizeRoutes(e *gin.Engine) {
	resize := e.Group("/resize")
	{
		resize.POST("", h.resize)
	}
}

func (h *Handler) resize(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}
