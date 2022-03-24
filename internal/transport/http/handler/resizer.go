package handler

import (
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
	err := h.services.Resizer.ResizeFromUrl(
		"https://cdn.pixabay.com/photo/2015/04/23/22/00/tree-736885__480.jpg",
		core.ResizeInput{
			Width:  500,
			Height: 100,
		},
	)
	if err != nil {
		h.setUnprocessableEntityJSONResponse(ctx, err.Error())
		return
	}
	// todo download the resized image

	h.setOkJSONResponse(ctx)
}
