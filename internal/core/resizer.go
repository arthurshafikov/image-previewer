package core

import "fmt"

type ResizeInput struct {
	Width  uint
	Height uint
}

type Image struct {
	Name      string
	Extension string
}

func (i *Image) GetFullName() string {
	return fmt.Sprintf("%s.%s", i.Name, i.Extension)
}
