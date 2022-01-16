package form3apiclient

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrRemoteError is a static error wrapped by all errors related to
// the remote server returning an error response.
var ErrRemoteError = errors.New("remote server returned an error")

// RemoteError constructs an error for a given HTTP status code.
func RemoteError(httpStatusCode int) error {
	return fmt.Errorf(
		"%w: http status code \"%d: %s\"",
		ErrRemoteError,
		httpStatusCode,
		http.StatusText(httpStatusCode))
}

// RemoteErrorWithServerMessage constructs an error for a given HTTP status code
// and an additional error message returned by the server.
func RemoteErrorWithServerMessage(httpStatusCode int, serverMessage string) error {
	return fmt.Errorf(
		"%w: http status code \"%d: %s\", server message: \"%s\"",
		ErrRemoteError,
		httpStatusCode,
		http.StatusText(httpStatusCode),
		serverMessage)
}

// ErrURLError is a static error wrapped by all errors related to
// problems with URL parsing.
var ErrURLError = errors.New("invalid url")

// URLError constructs an error for a given error message.
func URLError(message string) error {
	return fmt.Errorf("%w: %s", ErrURLError, message)
}

// WrapError wraps an external error and decorates it with an additional message.
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("error while %s: %w", message, err)
}
