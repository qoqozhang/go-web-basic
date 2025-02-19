package utils

import "time"

type MiddlewareJwt struct {
	Key              []byte
	SigningAlgorithm string
	MaxRefresh       time.Duration
}

func (m *MiddlewareJwt) Init() {

}
