package core

import "errors"

var (
	ErrServerError          = errors.New("500 server error")
	ErrWrongURL             = errors.New("given url is wrong")
	ErrOnlyJpg              = errors.New("only jpg and jpeg images are accepted")
	ErrCouldntDownloadImage = errors.New("couldn't download image from remote host")
	ErrCouldntSaveImage     = errors.New("couldn't save image to the local storage")
	ErrCouldntDecodeImage   = errors.New("couldn't decode image from jpeg")
)
