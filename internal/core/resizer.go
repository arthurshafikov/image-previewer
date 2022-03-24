package core

import (
	"fmt"
	"net/http"
	"os"
)

type ResizeInput struct {
	Header   http.Header
	ImageUrl string
	Width    int
	Height   int
}

type Image struct {
	Name      string
	Extension string
	File      *os.File
}

func (i *Image) GetFullName() string {
	return fmt.Sprintf("%s.%s", i.Name, i.Extension)
}
