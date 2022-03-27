package services

import (
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/oliamb/cutter"
	"github.com/thewolf27/image-previewer/internal/core"
)

type ResizerService struct {
	rawImageCache     ImageCache
	resizedImageCache ImageCache
}

func NewResizerService(rawImageCache ImageCache, resizedImageCache ImageCache) *ResizerService {
	return &ResizerService{
		rawImageCache:     rawImageCache,
		resizedImageCache: resizedImageCache,
	}
}

func (rs *ResizerService) ResizeFromUrl(inp core.ResizeInput) (*os.File, error) {
	resizedImage, err := rs.resizedImageCache.Remember(
		fmt.Sprintf("%s_%vx%v", inp.ImageUrl, inp.Width, inp.Height),
		func() (*core.Image, error) {
			image, err := rs.rawImageCache.Remember(inp.ImageUrl, func() (*core.Image, error) {
				return rs.downloadFromUrlAndSaveImageToStorage(inp)
			})
			if err != nil {
				return nil, err
			}

			resizedThumbnail, err := cutter.Crop(image.DecodedImage, cutter.Config{
				Width:  inp.Width,
				Height: inp.Height,
				Mode:   cutter.Centered,
			})
			if err != nil {
				return nil, err
			}

			resizedFile, err := os.Create(fmt.Sprintf(
				"%s/%s_%vx%v.%s",
				rs.resizedImageCache.GetCachedImagesFolder(),
				image.Name,
				inp.Width,
				inp.Height,
				image.Extension,
			))
			if err != nil {
				return nil, err
			}
			defer resizedFile.Close()

			if err := jpeg.Encode(resizedFile, resizedThumbnail, nil); err != nil {
				return nil, err
			}

			return &core.Image{
				File:         resizedFile,
				DecodedImage: resizedThumbnail,
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return resizedImage.File, nil
}

func (rs *ResizerService) downloadFromUrlAndSaveImageToStorage(inp core.ResizeInput) (*core.Image, error) {
	image, err := rs.parseImageNameFromUrl(inp.ImageUrl)
	if err != nil {
		return nil, err
	}

	body, err := rs.downloadImageFromUrl(inp.ImageUrl, inp.Header)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	image.File, err = rs.saveImageToStorage(image.GetFullName(), body)
	if err != nil {
		return nil, err
	}

	image.File.Seek(0, 0) // to avoid bug
	image.DecodedImage, err = jpeg.Decode(image.File)
	if err != nil {
		return nil, err
	}
	image.File.Close()

	return image, nil
}

func (rs *ResizerService) downloadImageFromUrl(url string, header http.Header) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (rs *ResizerService) saveImageToStorage(imageName string, body io.ReadCloser) (*os.File, error) {
	rawImageFile, err := os.Create(fmt.Sprintf(
		"%s/%s",
		rs.rawImageCache.GetCachedImagesFolder(),
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
