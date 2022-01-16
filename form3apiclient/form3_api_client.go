package form3apiclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
)

type form3APIRemoteError struct {
	ErrorMessage string `json:"error_message"`
}

type Form3ApiClient interface {
	Accounts() *accounts
}

type form3ApiClient struct {
	AccountsEndpoint *accounts
}

var config = restresourcehandler.Config{
	ResourceEncoding:     "application/json; charset=utf-8",
	IsDataWrapped:        true,
	DataPropertyName:     "data",
	RemoteErrorExtractor: extractRemoteError,
}

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

func NewForm3APIClient(apiURL string, httpClient *http.Client) *form3ApiClient {
	accounts, err := newAccounts(apiURL, httpClient)
	if err != nil {
		panic(err)
	}

	return &form3ApiClient{AccountsEndpoint: accounts}
}

func (c *form3ApiClient) Accounts() *accounts {
	return c.AccountsEndpoint
}
