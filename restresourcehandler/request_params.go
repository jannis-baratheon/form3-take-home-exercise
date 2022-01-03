package restresourcehandler

import (
	"fmt"
	"net/http"
)

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

func validateRequestParameters(params requestParams) {
	if !params.DoDiscardResourceId && params.ResourceId == "" {
		panic("Invalid request parameters: ResourceId is empty, but DoDiscardResourceId is not set.")
	}

	if !params.DoDiscardContent && params.Response == nil {
		panic("Invalid request parameters: Response is null, but DoDiscardContent is not set.")
	}

	switch params.HttpMethod {
	case http.MethodGet, http.MethodDelete, http.MethodPost:
	default:
		panic(fmt.Sprintf(`Unknown HTTP method "%s".`, params.HttpMethod))
	}
}
