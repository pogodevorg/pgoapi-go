package api

import (
	"errors"
	"fmt"

	protos "github.com/pogodevorg/pogo-protos"
)

// ErrFormatting happens when the something in the request body was not right
var ErrFormatting = errors.New("Request was malformatted and could not be performed")


// ErrRequest happens when there is an unknown issue with the request
var ErrRequest = errors.New("The request could not be completed")

// ErrNewRPCURL happens when there is a new RPC url available in the response body
var ErrNewRPCURL = errors.New("The request could not be completed")

// ErrInvalidAuthToken happens when the currently used auth token is not vailid
var ErrInvalidAuthToken = errors.New("The auth token used is not vailid")

// ErrRedirect happens when an invalid session endpoint has been used
var ErrRedirect = errors.New("The request was redirected")

// GetErrorFromStatus will, depending on the status code, give you an error or nil if there is no error
func GetErrorFromStatus(status protos.ResponseEnvelope_StatusCode) error {
	switch status {
	case protos.ResponseEnvelope_OK:
		return nil
	case protos.ResponseEnvelope_OK_RPC_URL_IN_RESPONSE:
		return ErrNewRPCURL
	case protos.ResponseEnvelope_REDIRECT:
		return ErrRedirect
	case protos.ResponseEnvelope_INVALID_AUTH_TOKEN:
		return ErrInvalidAuthToken
	default:
		return ErrRequest
	}
}

// ErrResponse happens when there's something wrong with the response object
type ErrResponse struct {
	err error
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("The response could not be read: %s", e.err.Error())
}
