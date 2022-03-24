package core

import (
	"fmt"
	"os"
)

type ResizeInput struct {
	Width  uint
	Height uint
}

type Image struct {
	Name      string
	Extension string
	File      *os.File
}

func (i *Image) GetFullName() string {
	return fmt.Sprintf("%s.%s", i.Name, i.Extension)
}
