package restresourcehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

func defaultRemoteErrorExtractor(response *http.Response) error {
	return RemoteError(response.StatusCode)
}

func createRequest(
	ctx context.Context,
	config Config,
	resourceURL url.URL,
	method string,
	id *string,
	queryParams map[string]string,
	resource interface{}) (*http.Request, error) {
	// append id
	if id != nil {
		resourceURL.Path = path.Join(resourceURL.Path, *id)
	}

	if queryParams != nil {
		query := resourceURL.Query()
		for key, val := range queryParams {
			query.Add(key, val)
		}

		resourceURL.RawQuery = query.Encode()
	}

	var body io.Reader

	if resource != nil {
		var err error
		body, err = readerForResource(config, resource)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, resourceURL.String(), body)
	if err != nil {
		return nil, WrapError(err, "constructing request")
	}

	return req, nil
}

func readResponse(config Config, reader io.Reader, response interface{}) error {
	respPayload, err := ioutil.ReadAll(reader)
	if err != nil {
		return WrapError(err, "decoding response")
	}

	if !config.IsDataWrapped {
		err = json.Unmarshal(respPayload, &response)

		return WrapError(err, "parsing response json")
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respPayload, &responseMap)

	if err != nil {
		return WrapError(err, "parsing response json")
	}

	err = json.Unmarshal(responseMap[config.DataPropertyName], &response)

	return WrapError(err, "parsing response json")
}

func readerForResource(config Config, resource interface{}) (io.Reader, error) {
	payload, err := json.Marshal(resource)
	if err != nil {
		return nil, WrapError(err, "marshalling request json")
	}

	if config.IsDataWrapped {
		payload, err = json.Marshal(map[string]json.RawMessage{"data": payload})
	}

	if err != nil {
		return nil, WrapError(err, "marshalling request json")
	}

	return bytes.NewReader(payload), nil
}
