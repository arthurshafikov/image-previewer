package core

import (
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/oliamb/cutter"
)

type DownloadImageInput struct {
	URL    string
	Header http.Header
}

type Image struct {
	Name         string
	Extension    string
	File         *os.File
	DecodedImage image.Image
}

func (i *Image) GetFullName() string {
	return fmt.Sprintf("%s.%s", i.Name, i.Extension)
}

func (i *Image) GetFullNameWithWidthAndHeight(width, height int) string {
	return fmt.Sprintf("%s_%vx%v.%s", i.Name, width, height, i.Extension)
}

func (i *Image) Crop(width, height int) (image.Image, error) {
	return cutter.Crop(i.DecodedImage, cutter.Config{
		Width:  width,
		Height: height,
		Mode:   cutter.Centered,
	})
}
