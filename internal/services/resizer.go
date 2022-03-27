package services

import (
	"fmt"
	"os"

	"github.com/oliamb/cutter"
	"github.com/thewolf27/image-previewer/internal/core"
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

func (rs *ResizerService) ResizeFromUrl(inp core.ResizeInput) (*os.File, error) {
	resizedImage, err := rs.resizedImageCache.Remember(
		fmt.Sprintf("%s_%vx%v", inp.ImageUrl, inp.Width, inp.Height),
		func() (*core.Image, error) {
			image, err := rs.rawImageCache.Remember(inp.ImageUrl, func() (*core.Image, error) {
				return rs.imagesService.DownloadFromUrlAndSaveImageToStorage(core.DownloadImageInput{
					Url:    inp.ImageUrl,
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
