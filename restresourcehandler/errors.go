package restresourcehandler

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrRemoteError = errors.New("remote server returned an error")

func RemoteError(httpStatusCode int) error {
	return fmt.Errorf(
		"%w: http status code \"%d: %s\"",
		ErrRemoteError,
		httpStatusCode,
		http.StatusText(httpStatusCode))
}

func RemoteErrorWithServerMessage(httpStatusCode int, serverMessage string) error {
	return fmt.Errorf(
		"%w: http status code \"%d: %s\", server message: \"%s\"",
		ErrRemoteError,
		httpStatusCode,
		http.StatusText(httpStatusCode),
		serverMessage)
}

func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("error while %s: %w", message, err)
}
