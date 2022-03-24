package services

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"github.com/thewolf27/image-previewer/internal/core"
)

const (
	rawImagesFolderName     = "raw"
	resizedImagesFolderName = "resized"
)

type ResizerService struct {
	storageFolder string
}

func NewResizerService(storageFolder string) *ResizerService {
	return &ResizerService{
		storageFolder: storageFolder,
	}
}

func (rs *ResizerService) ResizeFromUrl(url string, inp core.ResizeInput) error {
	image, err := rs.downloadAndSaveImageFromUrl(url)
	if err != nil {
		return err
	}

	image.File.Seek(0, 0) // to avoid bug
	decodedImageFile, err := jpeg.Decode(image.File)
	if err != nil {
		return err
	}
	image.File.Close()

	resizedThumbnail := resize.Thumbnail(inp.Width, inp.Height, decodedImageFile, resize.Lanczos3)

	resizedFile, err := os.Create(fmt.Sprintf(
		"%s/%s/%s_%vx%v.%s",
		rs.storageFolder,
		resizedImagesFolderName,
		image.Name,
		inp.Width,
		inp.Height,
		image.Extension,
	))
	if err != nil {
		return err
	}
	defer resizedFile.Close()

	if err := jpeg.Encode(resizedFile, resizedThumbnail, nil); err != nil {
		return err
	}

	return nil
}

func (rs *ResizerService) downloadAndSaveImageFromUrl(url string) (*core.Image, error) {
	image, err := rs.parseImageNameFromUrl(url)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	image.File, err = rs.saveImageToStorage(image.GetFullName(), resp.Body)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (rs *ResizerService) saveImageToStorage(imageName string, body io.ReadCloser) (*os.File, error) {
	rawImageFile, err := os.Create(fmt.Sprintf(
		"%s/%s/%s",
		rs.storageFolder,
		rawImagesFolderName,
		imageName,
	))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(rawImageFile, body); err != nil {
		return nil, err
	}

	return rawImageFile, nil
}

func (rs *ResizerService) parseImageNameFromUrl(url string) (*core.Image, error) {
	imageNameIndex := strings.LastIndex(url, "/")
	if imageNameIndex == -1 {
		return nil, core.ErrWrongUrl
	}

	fullImageName := url[imageNameIndex+1:]
	imageExtensionIndex := strings.LastIndex(fullImageName, ".")
	if imageExtensionIndex == -1 {
		return nil, core.ErrWrongUrl
	}

	return &core.Image{
		Name:      fullImageName[:imageExtensionIndex-1],
		Extension: fullImageName[imageExtensionIndex+1:],
	}, nil
}
