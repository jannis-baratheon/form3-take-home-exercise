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

type RestResourceHandler interface {
	Fetch(id string, params map[string]string, res interface{}) error
	Delete(id string, params map[string]string) error
	Create(resource interface{}, res interface{}) error
}

type restResourceHandler struct {
	HttpClient *http.Client
	Config     restResourceHandlerConfig
}

func NewRestResourceHandler(httpClient *http.Client, config restResourceHandlerConfig) RestResourceHandler {
	return &restResourceHandler{
		Config:     config,
		HttpClient: httpClient,
	}
}

// TODO parameterized
type errorDTO struct {
	ErrorMessage string `json:"error_message"`
}

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
	// copy base url
	u := c.Config.ResourceURL
	// append id
	if !params.DoDiscardResourceId {
		u.Path = path.Join(u.Path, params.ResourceId)
	}

	if params.QueryParams != nil {
		query := u.Query()
		for key, val := range params.QueryParams {
			query.Add(key, val)
		}
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if params.Resource != nil {
		payload, err := json.Marshal(params.Resource)

		if err != nil {
			return err
		}

		if c.Config.IsDataWrapped {
			payload, err = json.Marshal(map[string]json.RawMessage{"data": payload})
		}

		if err != nil {
			return err
		}

		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(params.HttpMethod, u.String(), body)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", c.Config.ResourceEncoding)
	req.Header.Add("Accept", c.Config.ResourceEncoding)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != params.ExpectedStatus {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var errRespJson errorDTO
		err = json.Unmarshal(body, &errRespJson)

		if err != nil {
			return err
		}

		return fmt.Errorf(`remote returned error status: %d, message: "%s"`, resp.StatusCode, errRespJson.ErrorMessage)
	}

	if params.DoDiscardContent {
		return nil
	}

	respPayload, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if !c.Config.IsDataWrapped {
		return json.Unmarshal(respPayload, &params.Response)
	}

	var responseMap map[string]json.RawMessage
	err = json.Unmarshal(respPayload, &responseMap)

	if err != nil {
		return err
	}

	return json.Unmarshal(responseMap[c.Config.DataPropertyName], &params.Response)
}
