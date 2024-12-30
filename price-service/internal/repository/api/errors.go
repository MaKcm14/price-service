package api

import "errors"

var (
	ErrServiceResponse  = errors.New("error of getting the response")
	ErrChromeDriver     = errors.New("error of the chrome driver's interaction")
	ErrBufferReading    = errors.New("error of reading the data")
	ErrConnectionClosed = errors.New("the client has closed the connection")
)
