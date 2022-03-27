package core

import "errors"

var (
	ErrWrongUrl = errors.New("given url is wrong")
	ErrOnlyJpg  = errors.New("only jpg and jpeg images are accepted")
)
