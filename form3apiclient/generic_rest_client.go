package form3apiclient

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
)

type GenericRestClient interface {
	Fetch(resourcePath string, id string, res interface{}) error
	// Delete(resourcePath string, id string) (*string, error)
	// Create(resourcePath string, json string) (*string, error)
}

type genericRestClientConfig struct {
	baseApiUrl url.URL
}

type genericRestClient struct {
	httpClient *http.Client
	config     genericRestClientConfig
}

func createGenericRestClient(config genericRestClientConfig) GenericRestClient {
	return &genericRestClient{
		httpClient: &http.Client{},
		config:     config,
	}
}

type errorDTO struct {
	ErrorMessage string `json:"error_message"`
}

type dataWrapperDTO struct {
	WrappedJSON json.RawMessage `json:"data"`
}

func (c *genericRestClient) Fetch(resourcePath string, id string, res interface{}) error {
	resourceUrl, err := c.config.baseApiUrl.Parse(fmt.Sprintf("%s/%s", resourcePath, id))

	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", resourceUrl.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		var errRespJson errorDTO
		err = json.Unmarshal(body, &errRespJson)

		if err != nil {
			return err
		}

		return fmt.Errorf(errRespJson.ErrorMessage)
	}

	var dataWrapper dataWrapperDTO
	err = json.Unmarshal(body, &dataWrapper)

	if err != nil {
		return err
	}

	err = json.Unmarshal(dataWrapper.WrappedJSON, &res)

	if err != nil {
		return err
	}

	return nil
}
