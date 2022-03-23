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

type ResizerService struct {
}

func NewResizerService() *ResizerService {
	return &ResizerService{}
}

func (rs *ResizerService) ResizeFromUrl(url string, inp core.ResizeInput) error {
	body, err := rs.downloadImageFromUrl(url)
	if err != nil {
		panic(err)
	}
	defer body.Close()

	image, err := rs.parseImageNameFromUrl(url)
	if err != nil {
		panic(err)
	}

	rawImageFile, err := rs.saveImageToStorage(image.GetFullName(), body)
	if err != nil {
		panic(err)
	}

	rawImageFile.Seek(0, 0)
	jpg, err := jpeg.Decode(rawImageFile)
	if err != nil {
		panic(err)
	}
	rawImageFile.Close()

	resized := resize.Resize(inp.Width, inp.Height, jpg, resize.Lanczos3)

	resizedFile, err := os.Create(fmt.Sprintf("./storage/cutted/%s_%v_%v.%s", image.Name, inp.Width, inp.Height, image.Extension))
	if err != nil {
		panic(err)
	}
	defer resizedFile.Close()

	if err := jpeg.Encode(resizedFile, resized, nil); err != nil {
		panic(err)
	}

	return nil
}

func (rs *ResizerService) downloadImageFromUrl(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (rs *ResizerService) saveImageToStorage(imageName string, body io.ReadCloser) (*os.File, error) {
	rawImageFile, err := os.Create(fmt.Sprintf("./storage/raw/%s", imageName))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(rawImageFile, body); err != nil {
		return nil, err
	}

	return rawImageFile, nil
}

func (rs *ResizerService) parseImageNameFromUrl(url string) (*core.Image, error) {
	index := strings.LastIndex(url, "/")
	if index == -1 {
		return nil, core.ErrWrongUrl
	}

	fullImageName := url[index+1:]
	extensionIndex := strings.LastIndex(fullImageName, ".")

	return &core.Image{
		Name:      fullImageName[:extensionIndex-1],
		Extension: fullImageName[extensionIndex+1:],
	}, nil
}
