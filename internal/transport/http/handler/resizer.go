package handler

import (
	"github.com/arthurshafikov/image-previewer/internal/core"
	"github.com/gin-gonic/gin"
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

	file, err := h.services.Resizer.ResizeFromUrl(
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

	ctx.Status(200)
	ctx.File(file.Name())
}
