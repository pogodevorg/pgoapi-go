package ptc

import (
	"fmt"
)

// LoginError is thrown when something went wrong with the login request
type LoginError struct {
	message string
}

func (e *LoginError) Error() string {
	return fmt.Sprintf("auth/ptc: %s", e.message)
}

func loginError(message string) (string, error) {
	return "", &LoginError{message}
}
