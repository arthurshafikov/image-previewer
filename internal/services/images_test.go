package services

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/arthurshafikov/image-previewer/internal/core"
	mock_services "github.com/arthurshafikov/image-previewer/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

var (
	someImageUrl = "https://some-url.com/some-image.jpg"
	coreImage    = &core.Image{
		Name:      "some-image",
		Extension: "jpg",
	}
)

func getImageService(t *testing.T) (*ImagesService, *mock_services.MockImageCache, *mock_services.MockImageCache) {
	t.Helper()

	ctl := gomock.NewController(t)
	rawImageCacheMock := mock_services.NewMockImageCache(ctl)
	resizedImageCacheMock := mock_services.NewMockImageCache(ctl)

	return NewImagesService(rawImageCacheMock, resizedImageCacheMock), rawImageCacheMock, resizedImageCacheMock
}

func TestDownloadImageFromUrl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	imagesService, _, _ := getImageService(t)

	httpmock.RegisterResponder(
		"GET",
		someImageUrl,
		httpmock.NewStringResponder(200, "123"),
	)

	resultBody, err := imagesService.downloadImageFromUrl(core.DownloadImageInput{
		Url:    someImageUrl,
		Header: http.Header{},
	})
	require.NoError(t, err)
	defer resultBody.Close()

	result, err := ioutil.ReadAll(resultBody)
	require.NoError(t, err)
	require.Equal(t, "123", string(result))
}

func TestParseImageNameFromUrl(t *testing.T) {
	imagesService, _, _ := getImageService(t)

	image, err := imagesService.parseImageNameFromUrl(someImageUrl)
	require.NoError(t, err)
	require.Equal(t, coreImage, image)
}

func TestParseImageNameFromUrlNotAnImage(t *testing.T) {
	imagesService, _, _ := getImageService(t)

	someNotImageUrl := "https://some-url.com/some.exe"

	image, err := imagesService.parseImageNameFromUrl(someNotImageUrl)
	require.ErrorIs(t, core.ErrOnlyJpg, err)
	require.Nil(t, image)
}

func TestParseImageNameFromUrlWrongUrl(t *testing.T) {
	imagesService, _, _ := getImageService(t)

	someNotImageUrl := "https://some-url.com/somePosts/2"

	image, err := imagesService.parseImageNameFromUrl(someNotImageUrl)
	require.ErrorIs(t, core.ErrWrongUrl, err)
	require.Nil(t, image)
}

func TestSaveRawImageToStorage(t *testing.T) {
	imagesService, rawImageCacheMock, _ := getImageService(t)
	gomock.InOrder(
		rawImageCacheMock.EXPECT().GetCachedImagesFolder().Times(1).Return("."),
	)
	body := io.NopCloser(strings.NewReader("someContent"))

	result, err := imagesService.saveRawImageToStorage("some-image.jpg", body)
	require.NoError(t, err)
	require.NotNil(t, result)

	fileContent, err := ioutil.ReadAll(result)
	require.NoError(t, err)
	require.Equal(t, "someContent", string(fileContent))
	require.NoError(t, os.Remove(result.Name()))
}

func TestSaveResizedImageToStorage(t *testing.T) {
	imagesService, _, resizedImageCacheMock := getImageService(t)
	gomock.InOrder(
		resizedImageCacheMock.EXPECT().GetCachedImagesFolder().Times(1).Return("."),
	)
	mockedImage := mock_services.NewImageMock()

	result, err := imagesService.SaveResizedImageToStorage(
		coreImage.GetFullNameWithWidthAndHeight(200, 500),
		mockedImage,
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "./some-image_200x500.jpg", result.Name())
	require.NoError(t, os.Remove(result.Name()))
}
