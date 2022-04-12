package core

import (
	"net/http"
)

type ResizeInput struct {
	Header   http.Header
	ImageURL string
	Width    int
	Height   int
}
