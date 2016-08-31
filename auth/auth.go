package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/pogodevorg/pgoapi-go/auth/google"
	"github.com/pogodevorg/pgoapi-go/auth/ptc"
)

// Provider is a common interface for managing auth tokens with the different third party authenticators
type Provider interface {
	Login(context.Context) (authToken string, err error)
	GetProviderString() string
	GetAccessToken() string
}

// UnknownProvider is a null provider for when a real one cannot be retrieved
type UnknownProvider struct {
}

// Login tries to log in
func (u *UnknownProvider) Login(ctx context.Context) (string, error) {
	return "", errors.New("Cannot log in using unknown provider")
}

// GetProviderString will return an identifying string for itself
func (u *UnknownProvider) GetProviderString() string {
	return "unknown"
}

// GetAccessToken will return an empty access token
func (u *UnknownProvider) GetAccessToken() string {
	return ""
}

// NewProvider creates a new provider based on the provider identifier
func NewProvider(provider, username, password string) (Provider, error) {
	switch provider {
	case "ptc":
		return ptc.NewProvider(username, password), nil
	case "google":
		return google.NewProvider(username, password), nil
	default:
		return &UnknownProvider{}, fmt.Errorf("Provider \"%s\" is not supported", provider)
	}
}
