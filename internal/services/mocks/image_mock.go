package mock_services //nolint

import (
	image "image"
	"image/color"
)

type ImageMock struct{}

func NewImageMock() *ImageMock {
	return &ImageMock{}
}

func (i *ImageMock) ColorModel() color.Model {
	return color.AlphaModel
}

func (i *ImageMock) Bounds() image.Rectangle {
	return image.Rect(0, 0, 0, 0)
}

func (i *ImageMock) At(x int, y int) color.Color {
	return color.Black
}
