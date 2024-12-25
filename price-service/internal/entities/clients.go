package entities

type ClientType int

const (
	APIClient ClientType = iota
	UserServiceClient
)
