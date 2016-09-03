package api

import (
	"errors"
	"fmt"

	protos "github.com/pogodevorg/POGOProtos-go"
)

// ErrFormatting happens when the something in the request body was not right
var ErrFormatting = errors.New("Request was malformatted and could not be performed")

// ErrNewRPCURL happens when there is a new RPC url available in the response body
var ErrNewRPCURL = errors.New("There is a new RPC url in the response, please use the new URL for future requests")

// ErrBadRequest happens when the remote service thinks the request is malformatted
var ErrBadRequest = errors.New("The remote service responded but appear to think the request is malformatted")

// ErrInvalidRequest happens when the request is invalid
var ErrInvalidRequest = errors.New("The remote service responded but appear to think the request is invalid")

// ErrInvalidPlatformRequest happens when a platform specific request like the request signature being incorrect
var ErrInvalidPlatformRequest = errors.New("A platform specific request is invalid")

// ErrSessionInvalidated happens when the session has been invalidated by the remote service
var ErrSessionInvalidated = errors.New("The session has been invalidated")

// ErrInvalidAuthToken happens when the currently used auth token is not vailid
var ErrInvalidAuthToken = errors.New("The auth token used is not vailid")

// ErrRedirect happens when an invalid session endpoint has been used
var ErrRedirect = errors.New("The request was redirected")

// ErrRequest happens when there is an unknown issue with the request
var ErrRequest = errors.New("The remote service responded but the request could not be completed for unknown reasons")

// ErrNoURL happens when the remote service is expected to respond with a remote URL but doesn't
var ErrNoURL = errors.New("The remote service did not respond with a remote URL when expected")

// ErrNoMapObjectsResponse happens when session.Announce is missing the map objects response
var ErrNoMapObjectsResponse = errors.New("The map objects response is missing")

// GetErrorFromStatus will, depending on the status code, give you an error or nil if there is no error
func GetErrorFromStatus(status protos.ResponseEnvelope_StatusCode) error {
	switch status {
	case protos.ResponseEnvelope_OK:
		return nil
	case protos.ResponseEnvelope_OK_RPC_URL_IN_RESPONSE:
		return ErrNewRPCURL
	case protos.ResponseEnvelope_BAD_REQUEST:
		return ErrBadRequest
	case protos.ResponseEnvelope_INVALID_REQUEST:
		return ErrInvalidRequest
	case protos.ResponseEnvelope_INVALID_PLATFORM_REQUEST:
		return ErrInvalidPlatformRequest
	case protos.ResponseEnvelope_REDIRECT:
		return ErrRedirect
	case protos.ResponseEnvelope_SESSION_INVALIDATED:
		return ErrSessionInvalidated
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
