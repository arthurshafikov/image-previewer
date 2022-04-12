package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arthurshafikov/image-previewer/internal/core"
	"github.com/arthurshafikov/image-previewer/internal/services"
	mock_services "github.com/arthurshafikov/image-previewer/internal/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResize(t *testing.T) {
	ctl := gomock.NewController(t)
	resizerServiceMock := mock_services.NewMockResizer(ctl)
	services := &services.Services{
		Resizer: resizerServiceMock,
	}
	w, c, h := getWriterContextAndHandler(t, services)
	expectedInput := core.ResizeInput{
		Header:   http.Header{"Someheaderkey": {"someHeaderValue"}},
		ImageURL: "https://some-website.com/some-image.jpg",
		Width:    200,
		Height:   500,
	}
	file, err := os.Create("some.jpg")
	require.NoError(t, err)
	defer require.NoError(t, file.Close())
	gomock.InOrder(
		resizerServiceMock.EXPECT().ResizeFromURL(expectedInput).Times(1).Return(file, nil),
	)
	c.Request = httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/resize/%v/%v/%s", expectedInput.Width, expectedInput.Height, expectedInput.ImageURL),
		nil,
	)
	c.Request.Header.Set("Someheaderkey", "someHeaderValue")
	c.Params = []gin.Param{
		{
			Key:   "width",
			Value: "200",
		},
		{
			Key:   "height",
			Value: "500",
		},
		{
			Key:   "imageURL",
			Value: "/https://some-website.com/some-image.jpg",
		},
	}

	h.resize(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, os.Remove("some.jpg"))
}

func TestResizeMissingParams(t *testing.T) {
	w, c, h := getWriterContextAndHandler(t, &services.Services{})
	c.Request = httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/resize/%v/%v/%s", 200, 500, "someURL"),
		nil,
	)
	c.Params = []gin.Param{}

	h.resize(c)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func getWriterContextAndHandler(
	t *testing.T,
	services *services.Services,
) (*httptest.ResponseRecorder, *gin.Context, *Handler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	h := NewHandler(context.Background(), services)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return w, c, h
}
