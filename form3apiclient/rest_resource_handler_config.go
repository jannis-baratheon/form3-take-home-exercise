package form3apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
)

// getRestResourceHandlerConfig constructs a configuration object for all
// Rest Resource Handlers used in this package.
func getRestResourceHandlerConfig() restresourcehandler.Config {
	return restresourcehandler.Config{
		ResourceEncoding:     "application/json; charset=utf-8",
		IsDataWrapped:        true,
		DataPropertyName:     "data",
		RemoteErrorExtractor: extractRemoteError,
	}
}

// extractRemoteError extracts additional information from the JSON sent along an error response.
func extractRemoteError(response *http.Response) error {
	if response.ContentLength == 0 {
		return RemoteError(response.StatusCode)
	}

	respPayload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WrapError(err, "reading response")
	}

	var remoteError form3APIRemoteError
	err = json.Unmarshal(respPayload, &remoteError)

	if err != nil {
		return WrapError(err, "parsing response json")
	}

	return RemoteErrorWithServerMessage(response.StatusCode, remoteError.ErrorMessage)
}
