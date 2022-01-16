package restresourcehandler

import (
	"context"
	"net/http"
	"net/url"
)

type RemoteErrorExtractor func(response *http.Response) error

type RestResourceHandler struct {
	client      *http.Client
	config      Config
	resourceURL url.URL
}

func NewRestResourceHandler(httpClient *http.Client, resourceURL string, config Config) *RestResourceHandler {
	validateRestResourceHandlerConfig(config)

	url, err := url.Parse(resourceURL)
	if err != nil {
		panic(err)
	}

	if !url.IsAbs() {
		panic("resource url must be absolute")
	}

	handler := RestResourceHandler{
		config:      config,
		client:      httpClient,
		resourceURL: *url,
	}

	return &handler
}

func (c *RestResourceHandler) Fetch(ctx context.Context, id string, params map[string]string, resp interface{}) error {
	return c.request(
		ctx,
		requestParams{
			HTTPMethod:     http.MethodGet,
			ResourceID:     id,
			QueryParams:    params,
			Response:       resp,
			ExpectedStatus: http.StatusOK,
		})
}

func (c *RestResourceHandler) Delete(ctx context.Context, id string, params map[string]string) error {
	return c.request(
		ctx,
		requestParams{
			HTTPMethod:       http.MethodDelete,
			ResourceID:       id,
			QueryParams:      params,
			DoDiscardContent: true,
			ExpectedStatus:   http.StatusNoContent,
		})
}

func (c *RestResourceHandler) Create(ctx context.Context, resourceToCreate interface{}, resp interface{}) error {
	return c.request(
		ctx,
		requestParams{
			HTTPMethod:          http.MethodPost,
			DoDiscardResourceID: true,
			Resource:            resourceToCreate,
			Response:            resp,
			ExpectedStatus:      http.StatusCreated,
		})
}

func (c *RestResourceHandler) request(ctx context.Context, params requestParams) error {
	validateRequestParameters(params)

	var id *string
	if !params.DoDiscardResourceID {
		id = &params.ResourceID
	}

	req, err := createRequest(ctx, c.config, c.resourceURL, params.HTTPMethod, id, params.QueryParams, params.Resource)
	if err != nil {
		return err
	}

	if !params.DoDiscardContent {
		req.Header.Add("Accept", c.config.ResourceEncoding)
	}

	if params.Resource != nil {
		req.Header.Add("Content-Type", c.config.ResourceEncoding)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return WrapError(err, "executing http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != params.ExpectedStatus {
		if c.config.RemoteErrorExtractor == nil {
			return defaultRemoteErrorExtractor(resp)
		}

		return c.config.RemoteErrorExtractor(resp)
	}

	if params.DoDiscardContent {
		return nil
	}

	return readResponse(c.config, resp.Body, params.Response)
}
