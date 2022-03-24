package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thewolf27/image-previewer/internal/core"
)

func (h *Handler) initResizeRoutes(e *gin.Engine) {
	resize := e.Group("/resize/:width/:height/*imageUrl")
	{
		resize.POST("", h.resize)
	}
}

func (h *Handler) resize(ctx *gin.Context) {
	width, err := h.parseIntegerFromParam(ctx, "width")
	if err != nil {
		h.setUnprocessableEntityJSONResponse(ctx, err.Error())
		return
	}
	height, err := h.parseIntegerFromParam(ctx, "height")
	if err != nil {
		h.setUnprocessableEntityJSONResponse(ctx, err.Error())
		return
	}
	imageUrl := ctx.Param("imageUrl")

	err = h.services.Resizer.ResizeFromUrl(
		core.ResizeInput{
			Header:   ctx.Request.Header,
			ImageUrl: imageUrl[1:], // avoid first slash
			Width:    width,
			Height:   height,
		},
	)
	if err != nil {
		h.setUnprocessableEntityJSONResponse(ctx, err.Error())
		return
	}
	// todo download the resized image

	h.setOkJSONResponse(ctx)
}
