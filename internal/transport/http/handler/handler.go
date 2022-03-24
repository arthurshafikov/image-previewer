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

func (h *Handler) parseUnsignedIntegerFromParam(ctx *gin.Context, key string) (uint, error) {
	param := ctx.Param(key)
	if param == "" {
		return 0, fmt.Errorf("the param %s is missing", key)
	}
	paramUint64, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(paramUint64), err
}

func (h *Handler) setUnprocessableEntityJSONResponse(ctx *gin.Context, data string) {
	h.setJSONResponse(ctx, http.StatusUnprocessableEntity, data)
}

func (h *Handler) setOkJSONResponse(ctx *gin.Context) {
	h.setJSONResponse(ctx, http.StatusOK, "OK")
}

func (h *Handler) setJSONResponse(ctx *gin.Context, code int, data string) {
	ctx.JSON(code, core.ServerResponse{
		Data: data,
	})
}
