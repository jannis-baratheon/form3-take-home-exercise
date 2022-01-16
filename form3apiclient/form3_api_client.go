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

type Form3ApiClient struct {
	accountsEndpoint *accounts
}

func getRestResourceHandlerConfig() restresourcehandler.Config {
	return restresourcehandler.Config{
		ResourceEncoding:     "application/json; charset=utf-8",
		IsDataWrapped:        true,
		DataPropertyName:     "data",
		RemoteErrorExtractor: extractRemoteError,
	}
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

func NewForm3APIClient(apiURL string, httpClient *http.Client) *Form3ApiClient {
	accounts, err := newAccounts(apiURL, httpClient)
	if err != nil {
		panic(err)
	}

	return &Form3ApiClient{accountsEndpoint: accounts}
}

func (c *Form3ApiClient) Accounts() *accounts {
	return c.accountsEndpoint
}
