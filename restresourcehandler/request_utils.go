package restresourcehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func defaultRemoteErrorExtractor(response *http.Response) error {
	return fmt.Errorf(`remote server returned error status: %d`, response.StatusCode)
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
