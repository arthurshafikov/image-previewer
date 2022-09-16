package services

import (
	"fmt"
	"os"

	"github.com/arthurshafikov/image-previewer/internal/core"
)

type ResizerService struct {
	logger            Logger
	rawImageCache     ImageCache
	resizedImageCache ImageCache
	imagesService     Images
}

func NewResizerService(
	logger Logger,
	rawImageCache ImageCache,
	resizedImageCache ImageCache,
	imagesService Images,
) *ResizerService {
	return &ResizerService{
		logger:            logger,
		rawImageCache:     rawImageCache,
		resizedImageCache: resizedImageCache,
		imagesService:     imagesService,
	}
}

func (rs *ResizerService) ResizeFromURL(inp core.ResizeInput) (*os.File, error) {
	resizedImage, err := rs.resizedImageCache.Remember(
		fmt.Sprintf("%s_%vx%v", inp.ImageURL, inp.Width, inp.Height),
		func() (*core.Image, error) {
			return rs.downloadImageAndResize(inp)
		},
	)
	if err != nil {
		return nil, err
	}

	return resizedImage.File, nil
}

func (rs *ResizerService) downloadImageAndResize(inp core.ResizeInput) (*core.Image, error) {
	image, err := rs.rawImageCache.Remember(inp.ImageURL, func() (*core.Image, error) {
		return rs.imagesService.DownloadFromURLAndSaveImageToStorage(core.DownloadImageInput{
			URL:    inp.ImageURL,
			Header: inp.Header,
		})
	})
	if err != nil {
		return nil, err
	}

	resizedThumbnail, err := image.Crop(inp.Width, inp.Height)
	if err != nil {
		rs.logger.Error(fmt.Errorf("error while cropping raw image: %w", err))

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
}
