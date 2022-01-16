package restresourcehandler

import (
	"fmt"
	"net/http"
)

// requestParams represents parameters of a REST API endpoint call.
type requestParams struct {
	// HTTPMethod is "GET" (fetch resource), "DELETE" (delete resource) or "POST" (resource creation).
	HTTPMethod string
	// ExpectedStatus is the HTTP status which will be considered a success.
	ExpectedStatus int
	// DoDiscardContent denotes if we expect a DTO in response content (false - content expected).
	DoDiscardContent bool
	// DoDiscardResourceID denotes if this request targets a particular resource
	// (true - no resource id, e.g. for POST requests).
	DoDiscardResourceID bool
	// ResourceID is the id of the resource targeted by this request.
	ResourceID string
	// QueryParams are additional query params to send with this request.
	QueryParams map[string]string
	// Resource is an object to be JSON-serialized and sent in this request.
	Resource interface{}
	// Response is an object that will be filled with the JSON-deserialized response content.
	Response interface{}
}

// validateRequestParameters does a sanity check of a requestParams instance.
func validateRequestParameters(params requestParams) {
	if !params.DoDiscardResourceID && params.ResourceID == "" {
		panic("Invalid request parameters: ResourceID is empty, but DoDiscardResourceID is not set.")
	}

	if !params.DoDiscardContent && params.Response == nil {
		panic("Invalid request parameters: Response is null, but DoDiscardContent is not set.")
	}

	switch params.HTTPMethod {
	case http.MethodGet, http.MethodDelete, http.MethodPost:
	default:
		panic(fmt.Sprintf(`Unknown HTTP method "%s".`, params.HTTPMethod))
	}
}
