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
	Fetch(resourcePath string, id string, params map[string]string, res interface{}, dataPropertyName string) error
	Delete(resourcePath string, id string, params map[string]string) error
	// Create(resourcePath string, json string) (*string, error)
}

type restClient struct {
	httpClient *http.Client
	baseApiUrl url.URL
}

// TODO max response size
func CreateRestClient(baseApiUrl url.URL, httpClient *http.Client) RestClient {
	return &restClient{
		baseApiUrl: baseApiUrl,
		httpClient: httpClient,
	}
}

// TODO parameterized
type errorDTO struct {
	ErrorMessage string `json:"error_message"`
}

func (c *restClient) Fetch(resourcePath string, id string, params map[string]string, res interface{}, dataPropertyName string) error {
	content, err := c.request("GET", http.StatusOK, false, resourcePath, id, params)
	
	if err != nil {
		return err
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(*content, &responseMap)

	if err != nil {
		return err
	}

	err = json.Unmarshal(responseMap[dataPropertyName], &res)

	return err
}

func (c *restClient) Delete(resourcePath string, id string, params map[string]string) error {
	_, err := c.request("DELETE", http.StatusNoContent, true, resourcePath, id, params)

	return err
}

func (c *restClient) request(method string, expectedStatus int, discardContent bool, resourcePath string, id string, params map[string]string) (*[]byte, error) {
	// copy base url
	u := c.baseApiUrl
	// join base url and relative resource url
	u.Path = path.Join(u.Path, fmt.Sprintf("/%s/%s", resourcePath, id))

	query := u.Query()
	for key, val := range params {
		query.Add(key, val)
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errRespJson errorDTO
		err = json.Unmarshal(body, &errRespJson)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(errRespJson.ErrorMessage)
	}

	if discardContent {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		return nil, err
	}

	return &body, nil
}
