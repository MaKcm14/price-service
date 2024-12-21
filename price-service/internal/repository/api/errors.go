package api

import "errors"

var (
	ErrServiceResponse error = errors.New("error of getting the response")
	ErrChromeDriver    error = errors.New("error of the chrome driver's interaction")
	ErrBufferReading   error = errors.New("error of reading the data")
)
