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

func defaultRemoteErrorExtractor(response *http.Response) error {
	return fmt.Errorf(`remote server returned error status: %d`, response.StatusCode)
}

func createRequest(
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

	return http.NewRequest(method, resourceURL.String(), body)
}

func readResponse(config Config, reader io.Reader, response interface{}) error {
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

func readerForResource(config Config, resource interface{}) (io.Reader, error) {
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
