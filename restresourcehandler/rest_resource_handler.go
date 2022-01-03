package restresourcehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

type RestResourceHandler interface {
	Fetch(id string, params map[string]string, res interface{}) error
	Delete(id string, params map[string]string) error
	Create(resource interface{}, res interface{}) error
}

type restResourceHandler struct {
	HttpClient *http.Client
	Config     restResourceHandlerConfig
}

func NewRestResourceHandler(httpClient *http.Client, config restResourceHandlerConfig) RestResourceHandler {
	return &restResourceHandler{
		Config:     config,
		HttpClient: httpClient,
	}
}

// TODO parameterized
type errorDTO struct {
	ErrorMessage string `json:"error_message"`
}

func (c *restResourceHandler) Fetch(id string, params map[string]string, res interface{}) error {
	return c.request("GET", http.StatusOK, false, &id, nil, params, res)
}

func (c *restResourceHandler) Delete(id string, params map[string]string) error {
	return c.request("DELETE", http.StatusNoContent, true, &id, nil, params, nil)
}

func (c *restResourceHandler) Create(resource interface{}, res interface{}) error {
	return c.request("POST", http.StatusCreated, false, nil, resource, nil, res)
}

func (c *restResourceHandler) request(method string, expectedStatus int, discardContent bool, id *string, resource interface{}, params map[string]string, res interface{}) error {
	// copy base url
	u := c.Config.ResourceURL
	// append id
	if id != nil {
		u.Path = path.Join(u.Path, *id)
	}

	if params != nil {
		query := u.Query()
		for key, val := range params {
			query.Add(key, val)
		}
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if resource != nil {
		payload, err := json.Marshal(resource)

		if err != nil {
			return err
		}

		if c.Config.IsDataWrapped {
			payload, err = json.Marshal(map[string]json.RawMessage{"data": payload})
		}

		if err != nil {
			return err
		}

		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, u.String(), body)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", c.Config.ResourceEncoding)
	req.Header.Add("Accept", c.Config.ResourceEncoding)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var errRespJson errorDTO
		err = json.Unmarshal(body, &errRespJson)

		if err != nil {
			return err
		}

		return fmt.Errorf(`remote returned error status: %d, message: "%s"`, resp.StatusCode, errRespJson.ErrorMessage)
	}

	if discardContent {
		return nil
	}

	respPayload, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if !c.Config.IsDataWrapped {
		return json.Unmarshal(respPayload, &res)
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respPayload, &responseMap)

	if err != nil {
		return err
	}

	return json.Unmarshal(responseMap[c.Config.DataPropertyName], &res)
}
