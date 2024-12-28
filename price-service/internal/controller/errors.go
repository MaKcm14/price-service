package controller

import "errors"

var (
	ErrRequestInfo    = errors.New("the wrong request data was got")
	ErrRequestPath    = errors.New("try to request to unknown resource")
	ErrServerHandling = errors.New("the server couldn't handle the response")
)

// ResponseErr is the wrapper for the errors' response.
type ResponseErr struct {
	Error string `json:"error"`
}
