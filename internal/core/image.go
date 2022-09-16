package core

import (
	"fmt"
	"image"
	"net/http"
	"os"
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
