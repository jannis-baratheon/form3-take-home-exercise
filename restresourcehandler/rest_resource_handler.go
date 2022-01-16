package restresourcehandler

import (
	"context"
	"net/http"
	"net/url"
)

// RestResourceHandler is used to query or update a REST API resource.
// Avoid cerating instances of RestResourceHandler directly.
// Rather use the NewRestResourceHandler function.
type RestResourceHandler struct {
	client      *http.Client
	config      Config
	resourceURL url.URL
}

// NewRestResourceHandler creates a new RestResourceHandler instance
// for a given Config, HTTP client isntance and resource URL
// (e.g. "http://example.com/api/qualifier/resource").
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

// Fetch fetches a resource for a given id, query parameters.
// resp is an output parameter that the fetched object will be stored in.
// Context can be used to control asynchronous requests.
func (c *RestResourceHandler) Fetch(
	ctx context.Context,
	resourceID string,
	queryParams map[string]string,
	resp interface{}) error {
	return c.request(
		ctx,
		requestParams{
			HTTPMethod:     http.MethodGet,
			ResourceID:     resourceID,
			QueryParams:    queryParams,
			Response:       resp,
			ExpectedStatus: http.StatusOK,
		})
}

// Delete deletes a resource with a given id.
// Additional query parameters can be specified to be sent with the request.
// Context can be used to control asynchronous requests.
func (c *RestResourceHandler) Delete(
	ctx context.Context,
	resourceID string,
	queryParams map[string]string) error {
	return c.request(
		ctx,
		requestParams{
			HTTPMethod:       http.MethodDelete,
			ResourceID:       resourceID,
			QueryParams:      queryParams,
			DoDiscardContent: true,
			ExpectedStatus:   http.StatusNoContent,
		})
}

// Create creates a resource given in the resourceToCreate parameter
// and stores the response in the resp parameter.
// Context can be used to control asynchronous requests.
func (c *RestResourceHandler) Create(
	ctx context.Context,
	resourceToCreate interface{},
	resp interface{}) error {
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
