package api

import (
	"errors"
)

// Crypto is a common interface for allowing pluggable signature
type Crypto interface {
	CreateIV() []byte
	Encrypt(in []byte, iv []byte) ([]byte, error)
	Enabled() bool
}

// DefaultCrypto is a minimally implemented version of the Crypto interface
type DefaultCrypto struct{}

// Encrypt does nothing in the mock
func (c *DefaultCrypto) Encrypt(in []byte, iv []byte) ([]byte, error) {
	return nil, errors.New("Cannot encrypt using the mock")
}

// CreateIV creates a random initialization vector to be used with the Encrypt function
func (c *DefaultCrypto) CreateIV() []byte {
	iv := make([]byte, 32)
	return iv
}

// Enabled returns wether or not crypto is enabled
func (c *DefaultCrypto) Enabled() bool {
	return false
}
