package restresourcehandler

import (
	"fmt"
	"net/http"
)

type requestParams struct {
	HTTPMethod          string
	ExpectedStatus      int
	DoDiscardContent    bool
	DoDiscardResourceID bool
	ResourceID          string
	QueryParams         map[string]string
	Resource            interface{}
	Response            interface{}
}

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
