package controller

import "errors"

var (
	ErrRequestInfo    = errors.New("the wrong request data was got")
	ErrServerHandling = errors.New("the server couldn't handle the response")
)
