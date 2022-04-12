package services

import (
	"net/http"
	"testing"

	"github.com/arthurshafikov/image-previewer/internal/core"
	mock_services "github.com/arthurshafikov/image-previewer/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestResizeFromURL(t *testing.T) {
	ctl := gomock.NewController(t)
	rawImageCacheMock := mock_services.NewMockImageCache(ctl)
	resizedImageCacheMock := mock_services.NewMockImageCache(ctl)
	imagesServiceMock := mock_services.NewMockImages(ctl)
	resizerService := NewResizerService(rawImageCacheMock, resizedImageCacheMock, imagesServiceMock)

	gomock.InOrder(
		resizedImageCacheMock.EXPECT().
			Remember(
				"https://some-url.com/some-image.jpg_200x500",
				gomock.AssignableToTypeOf(func() (*core.Image, error) { return nil, nil }),
			).
			Times(1).
			Return(coreImage, nil),
	)

	_, err := resizerService.ResizeFromURL(core.ResizeInput{
		ImageURL: someImageURL,
		Width:    200,
		Height:   500,
		Header:   http.Header{},
	})
	require.NoError(t, err)
}
