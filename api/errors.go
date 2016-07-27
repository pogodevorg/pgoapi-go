package api

const (
	successfulStatus            int32 = 1
	requiresAuthorizationStatus int32 = 2
	requiresRPCEndpointStatus   int32 = 53
	sessionTokenInvalidStatus   int32 = 102
)

// GetErrorFromStatus will, depending on the status code, give you an error or nil if there is no error
func GetErrorFromStatus(status int32) error {
	switch status {
	case successfulStatus:
		return nil
	case requiresAuthorizationStatus:
		return &RequiresAuthorizationError{}
	case requiresRPCEndpointStatus:
		return &RequiresRPCEndpointError{}
	case sessionTokenInvalidStatus:
		return &InvalidSessionError{}
	default:
		return &RequestError{}
	}
}

// RequestError happens when there's an error that is not accounted for
type RequestError struct{}

func (e *RequestError) Error() string {
	return "The request could not be completed"
}

// RequiresAuthorizationError happens when the API wants you to re-authorize the profile
type RequiresAuthorizationError struct{}

func (e *RequiresAuthorizationError) Error() string {
	return "The profile needs to be authorized"
}

// RequiresRPCEndpointError happens when an invalid session endpoint has been used
type RequiresRPCEndpointError struct{}

func (e *RequiresRPCEndpointError) Error() string {
	return "The request needs to be to a valid RPC endpoint"
}

// InvalidSessionError happens when the session token has expired or been invalidated
type InvalidSessionError struct{}

func (e *InvalidSessionError) Error() string {
	return "The session token is invalid"
}
