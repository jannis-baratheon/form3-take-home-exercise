package restclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type RestClient interface {
	Fetch(resourcePath string, id string, res interface{}) error
	// Delete(resourcePath string, id string) (*string, error)
	// Create(resourcePath string, json string) (*string, error)
}

type restClient struct {
	httpClient *http.Client
	baseApiUrl url.URL
}

func CreateRestClient(baseApiUrl url.URL, httpClient *http.Client) RestClient {
	return &restClient{
		baseApiUrl: baseApiUrl,
		httpClient: httpClient,
	}
}

type errorDTO struct {
	ErrorMessage string `json:"error_message"`
}

type dataWrapperDTO struct {
	WrappedJSON json.RawMessage `json:"data"`
}

func (c *restClient) Fetch(resourcePath string, id string, res interface{}) error {
	// copy base url
	u := c.baseApiUrl
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
