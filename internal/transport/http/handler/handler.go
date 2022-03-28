package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thewolf27/image-previewer/internal/core"
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

func (h *Handler) parseIntegerFromParam(ctx *gin.Context, key string) (int, error) {
	param := ctx.Param(key)
	if param == "" {
		return 0, fmt.Errorf("the param %s is missing", key)
	}
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return paramInt, err
}

func (h *Handler) setUnprocessableEntityJSONResponse(ctx *gin.Context, data string) {
	h.setJSONResponse(ctx, http.StatusUnprocessableEntity, data)
}

func (h *Handler) setJSONResponse(ctx *gin.Context, code int, data string) {
	ctx.JSON(code, core.ServerResponse{
		Data: data,
	})
}
