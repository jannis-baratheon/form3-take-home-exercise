package restresourcehandler

import (
	"bytes"
	"encoding/json"
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

type requestParams struct {
	HttpMethod          string
	ExpectedStatus      int
	DoDiscardContent    bool
	DoDiscardResourceId bool
	ResourceId          string
	QueryParams         map[string]string
	Resource            interface{}
	Response            interface{}
}

func (c *restResourceHandler) Fetch(id string, params map[string]string, res interface{}) error {
	return c.request(requestParams{
		HttpMethod:     "GET",
		ResourceId:     id,
		QueryParams:    params,
		Response:       res,
		ExpectedStatus: http.StatusOK})
}

func (c *restResourceHandler) Delete(id string, params map[string]string) error {
	return c.request(requestParams{
		HttpMethod:       "DELETE",
		ResourceId:       id,
		QueryParams:      params,
		DoDiscardContent: true,
		ExpectedStatus:   http.StatusNoContent})
}

func (c *restResourceHandler) Create(resource interface{}, res interface{}) error {
	return c.request(requestParams{
		HttpMethod:          "POST",
		DoDiscardResourceId: true,
		Resource:            resource,
		Response:            res,
		ExpectedStatus:      http.StatusCreated})
}

func (c *restResourceHandler) request(params requestParams) error {
	var err error

	// copy base url
	u := c.Config.ResourceURL
	// append id
	if !params.DoDiscardResourceId {
		u.Path = path.Join(u.Path, params.ResourceId)
	}

	if params.QueryParams != nil {
		query := u.Query()
		for key, val := range params.QueryParams {
			query.Add(key, val)
		}
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if params.Resource != nil {
		body, err = readerForResource(c.Config, params.Resource)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest(params.HttpMethod, u.String(), body)

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

	if resp.StatusCode != params.ExpectedStatus {
		return c.Config.RemoteErrorExtractor(resp)
	}

	if params.DoDiscardContent {
		return nil
	}

	return readResponse(c.Config, resp.Body, params.Response)
}

func readResponse(config restResourceHandlerConfig, reader io.Reader, response interface{}) error {
	respPayload, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	if !config.IsDataWrapped {
		return json.Unmarshal(respPayload, &response)
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respPayload, &responseMap)

	if err != nil {
		return err
	}

	return json.Unmarshal(responseMap[config.DataPropertyName], &response)
}

func readerForResource(config restResourceHandlerConfig, resource interface{}) (io.Reader, error) {
	payload, err := json.Marshal(resource)

	if err != nil {
		return nil, err
	}

	if config.IsDataWrapped {
		payload, err = json.Marshal(map[string]json.RawMessage{"data": payload})
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(payload), nil
}
