package core

import (
	"net/http"
)

type ResizeInput struct {
	Header   http.Header
	ImageUrl string
	Width    int
	Height   int
}
