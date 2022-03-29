package services

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/thewolf27/image-previewer/internal/core"
)

type ImagesService struct {
	rawImageCache     ImageCache
	resizedImageCache ImageCache
}

func NewImagesService(rawImageCache ImageCache, resizedImageCache ImageCache) *ImagesService {
	return &ImagesService{
		rawImageCache:     rawImageCache,
		resizedImageCache: resizedImageCache,
	}
}

func (is *ImagesService) DownloadFromUrlAndSaveImageToStorage(inp core.DownloadImageInput) (*core.Image, error) {
	image, err := is.parseImageNameFromUrl(inp.Url)
	if err != nil {
		return nil, err
	}

	body, err := is.downloadImageFromUrl(inp)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	image.File, err = is.saveRawImageToStorage(image.GetFullName(), body)
	if err != nil {
		return nil, err
	}

	image.DecodedImage, err = jpeg.Decode(image.File)
	if err != nil {
		return nil, err
	}
	image.File.Close()

	return image, nil
}

func (is *ImagesService) SaveResizedImageToStorage(imageName string, resizedImage image.Image) (*os.File, error) {
	resizedFile, err := os.Create(fmt.Sprintf(
		"%s/%s",
		is.resizedImageCache.GetCachedImagesFolder(),
		imageName,
	))
	if err != nil {
		return nil, err
	}
	defer resizedFile.Close()

	if err := jpeg.Encode(resizedFile, resizedImage, nil); err != nil {
		return nil, err
	}

	return resizedFile, nil
}

func (is *ImagesService) parseImageNameFromUrl(url string) (*core.Image, error) {
	imageNameIndex := strings.LastIndex(url, "/")
	if imageNameIndex == -1 {
		return nil, core.ErrWrongUrl
	}

	fullImageName := url[imageNameIndex+1:]
	imageExtensionIndex := strings.LastIndex(fullImageName, ".")
	if imageExtensionIndex == -1 {
		return nil, core.ErrWrongUrl
	}
	imageExtension := fullImageName[imageExtensionIndex+1:]

	if imageExtension != "jpg" && imageExtension != "jpeg" {
		return nil, core.ErrOnlyJpg
	}

	return &core.Image{
		Name:      fullImageName[:imageExtensionIndex],
		Extension: imageExtension,
	}, nil
}

func (is *ImagesService) downloadImageFromUrl(inp core.DownloadImageInput) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", inp.Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = inp.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("remote host has returned: " + resp.Status)
	}

	return resp.Body, nil
}

func (is *ImagesService) saveRawImageToStorage(imageName string, body io.ReadCloser) (*os.File, error) {
	rawImageFile, err := os.Create(fmt.Sprintf(
		"%s/%s",
		is.rawImageCache.GetCachedImagesFolder(),
		imageName,
	))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(rawImageFile, body); err != nil {
		return nil, err
	}
	rawImageFile.Seek(0, 0) // to avoid bug

	return rawImageFile, nil
}
