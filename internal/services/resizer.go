package services

import (
	"fmt"
	"os"

	"github.com/arthurshafikov/image-previewer/internal/core"
	"github.com/oliamb/cutter"
)

type ResizerService struct {
	rawImageCache     ImageCache
	resizedImageCache ImageCache
	imagesService     Images
}

func NewResizerService(
	rawImageCache ImageCache,
	resizedImageCache ImageCache,
	imagesService Images,
) *ResizerService {
	return &ResizerService{
		rawImageCache:     rawImageCache,
		resizedImageCache: resizedImageCache,
		imagesService:     imagesService,
	}
}

func (rs *ResizerService) ResizeFromURL(inp core.ResizeInput) (*os.File, error) {
	resizedImage, err := rs.resizedImageCache.Remember(
		fmt.Sprintf("%s_%vx%v", inp.ImageURL, inp.Width, inp.Height),
		func() (*core.Image, error) {
			image, err := rs.rawImageCache.Remember(inp.ImageURL, func() (*core.Image, error) {
				return rs.imagesService.DownloadFromURLAndSaveImageToStorage(core.DownloadImageInput{
					URL:    inp.ImageURL,
					Header: inp.Header,
				})
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

			resizedFile, err := rs.imagesService.SaveResizedImageToStorage(
				image.GetFullNameWithWidthAndHeight(inp.Width, inp.Height),
				resizedThumbnail,
			)
			if err != nil {
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
