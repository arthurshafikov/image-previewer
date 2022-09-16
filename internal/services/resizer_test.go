package services

import (
	"image"
	"net/http"
	"os"
	"testing"

	"github.com/arthurshafikov/image-previewer/internal/core"
	mock_services "github.com/arthurshafikov/image-previewer/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var input = core.ResizeInput{
	ImageURL: someImageURL,
	Width:    200,
	Height:   500,
	Header:   http.Header{},
}

func getResizerService(t *testing.T) (
	*ResizerService,
	*mock_services.MockImageCache,
	*mock_services.MockImageCache,
	*mock_services.MockImages,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	rawImageCacheMock := mock_services.NewMockImageCache(ctrl)
	resizedImageCacheMock := mock_services.NewMockImageCache(ctrl)
	imagesServiceMock := mock_services.NewMockImages(ctrl)
	loggerMock := mock_services.NewMockLogger(ctrl)
	resizerService := NewResizerService(loggerMock, rawImageCacheMock, resizedImageCacheMock, imagesServiceMock)

	return resizerService, rawImageCacheMock, resizedImageCacheMock, imagesServiceMock
}

func TestResizeFromURL(t *testing.T) {
	resizerService, _, resizedImageCacheMock, _ := getResizerService(t)
	gomock.InOrder(
		resizedImageCacheMock.EXPECT().Remember(
			"https://some-url.com/some-image.jpg_200x500",
			gomock.AssignableToTypeOf(func() (*core.Image, error) { return nil, nil }),
		).Return(coreImage, nil),
	)

	_, err := resizerService.ResizeFromURL(input)
	require.NoError(t, err)
}

func TestDownloadImageAndResize(t *testing.T) {
	resizerService, rawImageCacheMock, _, imagesServiceMock := getResizerService(t)
	coreImage.DecodedImage = image.Black

	gomock.InOrder(
		rawImageCacheMock.EXPECT().Remember(
			"https://some-url.com/some-image.jpg",
			gomock.AssignableToTypeOf(func() (*core.Image, error) { return nil, nil }),
		).Return(coreImage, nil),
		imagesServiceMock.EXPECT().SaveResizedImageToStorage("some-image_200x500.jpg", gomock.Any()).Return(&os.File{}, nil),
	)

	result, err := resizerService.downloadImageAndResize(input)
	require.NoError(t, err)
	require.NotNil(t, result)
}
