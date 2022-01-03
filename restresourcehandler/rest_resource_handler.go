package restresourcehandler

import (
	"net/http"
)

type RemoteErrorExtractor func(response *http.Response) error

type RestResourceHandler interface {
	Fetch(id string, params map[string]string, res interface{}) error
	Delete(id string, params map[string]string) error
	Create(resource interface{}, res interface{}) error
}

type restResourceHandler struct {
	HttpClient *http.Client
	Config     RestResourceHandlerConfig
}

func NewRestResourceHandler(httpClient *http.Client, config RestResourceHandlerConfig) RestResourceHandler {
	validateRestResourceHandlerConfig(config)

	return &restResourceHandler{
		Config:     config,
		HttpClient: httpClient,
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
