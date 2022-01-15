package restresourcehandler

import (
	"net/http"
	"net/url"
)

type RemoteErrorExtractor func(response *http.Response) error

type RestResourceHandler interface {
	Fetch(id string, params map[string]string, response interface{}) error
	Delete(id string, params map[string]string) error
	Create(resourceToCreate interface{}, response interface{}) error
}

type restResourceHandler struct {
	HTTPClient  *http.Client
	Config      Config
	ResourceURL url.URL
}

func NewRestResourceHandler(httpClient *http.Client, resourceURL string, config Config) RestResourceHandler {
	validateRestResourceHandlerConfig(config)

	url, err := url.Parse(resourceURL)
	if err != nil {
		panic(err)
	}

	if !url.IsAbs() {
		panic("resource url must be absolute")
	}

	handler := restResourceHandler{
		Config:      config,
		HTTPClient:  httpClient,
		ResourceURL: *url,
	}

	return &handler
}

func (c *restResourceHandler) Fetch(id string, params map[string]string, response interface{}) error {
	return c.request(requestParams{
		HTTPMethod:     http.MethodGet,
		ResourceID:     id,
		QueryParams:    params,
		Response:       response,
		ExpectedStatus: http.StatusOK,
	})
}

func (c *restResourceHandler) Delete(id string, params map[string]string) error {
	return c.request(requestParams{
		HTTPMethod:       http.MethodDelete,
		ResourceID:       id,
		QueryParams:      params,
		DoDiscardContent: true,
		ExpectedStatus:   http.StatusNoContent,
	})
}

func (c *restResourceHandler) Create(resourceToCreate interface{}, response interface{}) error {
	return c.request(requestParams{
		HTTPMethod:          http.MethodPost,
		DoDiscardResourceID: true,
		Resource:            resourceToCreate,
		Response:            response,
		ExpectedStatus:      http.StatusCreated,
	})
}

func (c *restResourceHandler) request(params requestParams) error {
	validateRequestParameters(params)

	var id *string
	if !params.DoDiscardResourceID {
		id = &params.ResourceID
	}
	req, err := createRequest(c.Config, c.ResourceURL, params.HTTPMethod, id, params.QueryParams, params.Resource)
	if err != nil {
		return err
	}

	if !params.DoDiscardContent {
		req.Header.Add("Accept", c.Config.ResourceEncoding)
	}

	if params.Resource != nil {
		req.Header.Add("Content-Type", c.Config.ResourceEncoding)
	}

	resp, err := c.HTTPClient.Do(req)
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
