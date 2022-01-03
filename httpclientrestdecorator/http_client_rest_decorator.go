package httpclientrestdecorator

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

type GenericRestClient interface {
	Fetch(resourcePath string, id string, res interface{}) error
	// Delete(resourcePath string, id string) (*string, error)
	// Create(resourcePath string, json string) (*string, error)
}

type GenericRestClientConfig struct {
	BaseApiUrl url.URL
}

type genericRestClient struct {
	httpClient *http.Client
	config     GenericRestClientConfig
}

func CreateGenericRestClient(config GenericRestClientConfig) GenericRestClient {
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
	// copy base url
	u := c.config.BaseApiUrl
	// join base url and relative resource url
	u.Path = path.Join(u.Path, fmt.Sprintf("/%s/%s", resourcePath, id)) 

	req, err := http.NewRequest("GET", u.String(), nil)
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

	if resp.StatusCode != http.StatusOK {
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
