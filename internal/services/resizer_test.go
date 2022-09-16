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

func TestResizeFromURL(t *testing.T) {
	ctl := gomock.NewController(t)
	rawImageCacheMock := mock_services.NewMockImageCache(ctl)
	resizedImageCacheMock := mock_services.NewMockImageCache(ctl)
	imagesServiceMock := mock_services.NewMockImages(ctl)
	resizerService := NewResizerService(rawImageCacheMock, resizedImageCacheMock, imagesServiceMock)

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
	ctl := gomock.NewController(t)
	rawImageCacheMock := mock_services.NewMockImageCache(ctl)
	resizedImageCacheMock := mock_services.NewMockImageCache(ctl)
	imagesServiceMock := mock_services.NewMockImages(ctl)
	resizerService := NewResizerService(rawImageCacheMock, resizedImageCacheMock, imagesServiceMock)
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
