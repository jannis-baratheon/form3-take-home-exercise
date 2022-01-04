package form3apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/jannis-baratheon/Form3-take-home-excercise/restresourcehandler"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type form3APIRemoteError struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

type Form3ApiClient interface {
	CreateAccount(accountData AccountData) (AccountData, error)
	GetAccount(id string) (AccountData, error)
	DeleteAccount(id string, version int) error
}

type form3ApiClient struct {
	AccountHandler restresourcehandler.RestResourceHandler
}

var config = restresourcehandler.RestResourceHandlerConfig{
	ResourceEncoding:     "application/json; charset=utf-8",
	IsDataWrapped:        true,
	DataPropertyName:     "data",
	RemoteErrorExtractor: extractRemoteError,
}

func extractRemoteError(response *http.Response) error {
	respPayload, err := ioutil.ReadAll(response.Body)
	if response.ContentLength == 0 {
		return fmt.Errorf(`api responded with error: http status code %d, http status "%s"`, response.StatusCode, response.Status)
	}

	if err != nil {
		return err
	}

	var remoteError form3APIRemoteError
	err = json.Unmarshal(respPayload, &remoteError)

	if err != nil {
		return err
	}

	return fmt.Errorf(`api responded with error: http status code %d, http status "%s", server message: "%s"`, response.StatusCode, response.Status, remoteError.ErrorMessage)
}

func NewForm3APIClient(apiURL string, httpClient *http.Client) Form3ApiClient {
	url, err := url.Parse(apiURL)

	if err != nil {
		panic(err)
	}

	if !url.IsAbs() {
		panic(fmt.Errorf("api url must be absolute"))
	}

	url.Path = path.Join(url.Path, "organisation/accounts")

	accountHandler := restresourcehandler.NewRestResourceHandler(httpClient, url.String(), config)

	return &form3ApiClient{AccountHandler: accountHandler}
}

func (c *form3ApiClient) GetAccount(id string) (AccountData, error) {
	var accountData AccountData
	err := c.AccountHandler.Fetch(id, nil, &accountData)

	if err != nil {
		return accountData, err
	}

	return accountData, nil
}

func (c *form3ApiClient) DeleteAccount(id string, version int) error {
	return c.AccountHandler.Delete(id, map[string]string{"version": fmt.Sprint(version)})
}

func (c *form3ApiClient) CreateAccount(accountData AccountData) (AccountData, error) {
	var response AccountData
	err := c.AccountHandler.Create(&accountData, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}
