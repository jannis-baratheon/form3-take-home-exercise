package restresourcehandler

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrRemoteError is a static error wrapped by all errors related to
// the remote server returning an error response.
var ErrRemoteError = errors.New("remote server returned an error")

// RemoteError constructs an error for the given HTTP status code.
func RemoteError(httpStatusCode int) error {
	return fmt.Errorf(
		"%w: http status code \"%d: %s\"",
		ErrRemoteError,
		httpStatusCode,
		http.StatusText(httpStatusCode))
}

// WrapError wraps an external error and decorates it with an additional message.
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("error while %s: %w", message, err)
}
