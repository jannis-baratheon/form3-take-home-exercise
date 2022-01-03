package restresourcehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type RemoteErrorExtractor func(response *http.Response) error

type RestResourceHandler interface {
	Fetch(id string, params map[string]string, res interface{}) error
	Delete(id string, params map[string]string) error
	Create(resource interface{}, res interface{}) error
}

// TODO maxresponsesize, errordeserializer
type RestResourceHandlerConfig struct {
	RemoteErrorExtractor RemoteErrorExtractor
	ResourceURL          url.URL
	ResourceEncoding     string
	DataPropertyName     string
	IsDataWrapped        bool
}

type restResourceHandler struct {
	HttpClient *http.Client
	Config     RestResourceHandlerConfig
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

func NewRestResourceHandler(httpClient *http.Client, config RestResourceHandlerConfig) RestResourceHandler {
	validateRestResourceHandlerConfig(config)

	return &restResourceHandler{
		Config:     config,
		HttpClient: httpClient,
	}
}

func validateRestResourceHandlerConfig(config RestResourceHandlerConfig) {
	if config.IsDataWrapped && config.DataPropertyName == "" {
		panic("IsDataWrapped is set, but DataPropertyName has not been given.")
	}

	if !config.IsDataWrapped && config.DataPropertyName != "" {
		panic("IsDataWrapped is not set, but DataPropertyName has been given.")
	}

	if !config.ResourceURL.IsAbs() {
		panic("Resource URL must be absolute.")
	}

	if config.ResourceEncoding == "" {
		panic("ResourceEncoding must be set.")
	}
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
	validateRequestParameters(params)

	var id *string
	if !params.DoDiscardResourceId {
		id = &params.ResourceId
	}
	req, err := createRequest(c.Config, params.HttpMethod, id, params.QueryParams, params.Resource)

	if err != nil {
		return err
	}

	if !params.DoDiscardContent {
		req.Header.Add("Accept", c.Config.ResourceEncoding)
	}

	if params.Resource != nil {
		req.Header.Add("Content-Type", c.Config.ResourceEncoding)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != params.ExpectedStatus {
		if c.Config.RemoteErrorExtractor == nil {
			return defaultRemoteErrorExtractor(resp)
		}
		return c.Config.RemoteErrorExtractor(resp)
	}

	if params.DoDiscardContent {
		return nil
	}

	return readResponse(c.Config, resp.Body, params.Response)
}

func validateRequestParameters(params requestParams) {
	if !params.DoDiscardResourceId && params.ResourceId == "" {
		panic("Invalid request parameters: ResourceId is empty, but DoDiscardResourceId is not set.")
	}

	if !params.DoDiscardContent && params.Response == nil {
		panic("Invalid request parameters: Response is null, but DoDiscardContent is not set.")
	}

	switch params.HttpMethod {
	case "GET", "POST", "DELETE":
	default:
		panic(fmt.Sprintf(`Unknown HTTP method "%s".`, params.HttpMethod))
	}
}

func defaultRemoteErrorExtractor(response *http.Response) error {
	return fmt.Errorf(`remote server returned error status: %d"`, response.StatusCode)
}

func createRequest(config RestResourceHandlerConfig, method string, id *string, queryParams map[string]string, resource interface{}) (*http.Request, error) {
	var err error

	// copy base url
	u := config.ResourceURL
	// append id
	if id != nil {
		u.Path = path.Join(u.Path, *id)
	}

	if queryParams != nil {
		query := u.Query()
		for key, val := range queryParams {
			query.Add(key, val)
		}
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if resource != nil {
		body, err = readerForResource(config, resource)
		if err != nil {
			return nil, err
		}
	}

	return http.NewRequest(method, u.String(), body)
}

func readResponse(config RestResourceHandlerConfig, reader io.Reader, response interface{}) error {
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

func readerForResource(config RestResourceHandlerConfig, resource interface{}) (io.Reader, error) {
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
