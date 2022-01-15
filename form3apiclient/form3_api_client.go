package form3apiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
)

type form3APIRemoteError struct {
	ErrorMessage string `json:"error_message"`
}

type Form3ApiClient interface {
	Accounts() Accounts
}

type form3ApiClient struct {
	AccountsEndpoint Accounts
}

var config = restresourcehandler.Config{
	ResourceEncoding:     "application/json; charset=utf-8",
	IsDataWrapped:        true,
	DataPropertyName:     "data",
	RemoteErrorExtractor: extractRemoteError,
}

func extractRemoteError(response *http.Response) error {
	respPayload, err := ioutil.ReadAll(response.Body)

	if response.ContentLength == 0 {
		return fmt.Errorf(
			`api responded with error: http status code %d, http status "%s"`,
			response.StatusCode,
			response.Status)
	}

	if err != nil {
		return err
	}

	var remoteError form3APIRemoteError
	err = json.Unmarshal(respPayload, &remoteError)

	if err != nil {
		return err
	}

	return fmt.Errorf(
		`api responded with error: http status code %d, http status "%s", server message: "%s"`,
		response.StatusCode,
		response.Status,
		remoteError.ErrorMessage)
}

func NewForm3APIClient(apiURL string, httpClient *http.Client) Form3ApiClient {
	accounts, err := newAccounts(apiURL, httpClient)
	if err != nil {
		panic(err)
	}

	return &form3ApiClient{AccountsEndpoint: accounts}
}

func (c *form3ApiClient) Accounts() Accounts {
	return c.AccountsEndpoint
}
