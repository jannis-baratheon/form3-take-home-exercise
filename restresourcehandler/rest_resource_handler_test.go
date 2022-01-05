package restresourcehandler_test

import (
	"encoding/json"
	"fmt"
	"github.com/jannis-baratheon/Form3-take-home-excercise/restresourcehandler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"io/ioutil"
	"net/http"
)

type person struct {
	Name string `json:"name"`
}

type wrapper struct {
	Data person `json:"data"`
}

type apiError struct {
	ErrorMessage string `json:"error_message"`
}

type apiCall func(client restresourcehandler.RestResourceHandler) error

var exampleValidApiCalls = map[string]apiCall{
	"fetch": func(client restresourcehandler.RestResourceHandler) error {
		var response person
		return client.Fetch("1", map[string]string{"attrs": "name"}, &response)
	},
	"delete": func(client restresourcehandler.RestResourceHandler) error {
		return client.Delete("1", map[string]string{"version": "1"})
	},
	"create": func(client restresourcehandler.RestResourceHandler) error {
		var response person
		return client.Create(person{"Smith"}, &response)
	},
}

func forEachExampleValidApiCall(consumer func(string, apiCall)) {
	for callName, call := range exampleValidApiCalls {
		consumer(callName, call)
	}
}

var _ = Describe("RestResourceHandler", func() {
	var server *ghttp.Server
	var httpClient *http.Client
	var url string

	const resourceEncoding = "application/json; charset=utf-8"
	const resourcePath = "/api/people"

	BeforeEach(func() {
		server = ghttp.NewServer()
		url = server.URL() + resourcePath
		httpClient = &http.Client{}
	})

	AfterEach(func() {
		server.Close()
	})

	When("on happy-path", func() {
		var client restresourcehandler.RestResourceHandler

		BeforeEach(func() {
			client = restresourcehandler.NewRestResourceHandler(
				httpClient,
				url,
				restresourcehandler.RestResourceHandlerConfig{
					IsDataWrapped:    true,
					DataPropertyName: "data",
					ResourceEncoding: resourceEncoding,
				})
		})

		It("fetches resource", func() {
			expectedPerson := person{"Smith"}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", resourcePath+"/1", "attrs=name"),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.RespondWithJSONEncoded(http.StatusOK, wrapper{expectedPerson})))

			var response person
			err := client.Fetch("1", map[string]string{"attrs": "name"}, &response)

			Expect(err).To(Succeed())
			Expect(response).To(Equal(expectedPerson))
		})

		It("deletes resource", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", resourcePath+"/1", "version=1"),
					ghttp.RespondWith(http.StatusNoContent, nil)))

			err := client.Delete("1", map[string]string{"version": "1"})

			Expect(err).NotTo(HaveOccurred())
		})

		It("creates resource", func() {
			payload := person{"Smith"}
			expectedResponse := person{"Gennings"}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", resourcePath),
					ghttp.VerifyContentType(resourceEncoding),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.VerifyJSONRepresenting(wrapper{payload}),
					ghttp.RespondWithJSONEncoded(http.StatusCreated, wrapper{expectedResponse})))

			var actualResponse person
			err := client.Create(payload, &actualResponse)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualResponse).To(Equal(expectedResponse))
		})
	})

	Context("with default remote error extractor", func() {
		var client restresourcehandler.RestResourceHandler

		BeforeEach(func() {
			client = restresourcehandler.NewRestResourceHandler(
				httpClient,
				url,
				restresourcehandler.RestResourceHandlerConfig{
					ResourceEncoding: resourceEncoding,
				})
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil)))
		})

		forEachExampleValidApiCall(func(reqName string, req apiCall) {
			It(fmt.Sprintf(`reports default remote error during "%s" call`, reqName), func() {
				err := req(client)

				Expect(err).To(MatchError(fmt.Errorf("remote server returned error status: 500")))
			})
		})
	})

	Context("with custom remote error extractor returning an error not based on response", func() {
		var client restresourcehandler.RestResourceHandler
		customError := fmt.Errorf("some custom error")

		BeforeEach(func() {
			client = restresourcehandler.NewRestResourceHandler(
				httpClient,
				url,
				restresourcehandler.RestResourceHandlerConfig{
					ResourceEncoding: resourceEncoding,
					RemoteErrorExtractor: func(response *http.Response) error {
						return customError
					},
				})
		})

		forEachExampleValidApiCall(func(reqName string, req apiCall) {
			It(fmt.Sprintf(`reports custom remote error during "%s" call`, reqName), func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusInternalServerError, nil)))

				err := req(client)

				Expect(err).To(MatchError(customError))
			})
		})
	})

	Context("with custom remote error extractor returning error based on message from response", func() {
		var client restresourcehandler.RestResourceHandler

		BeforeEach(func() {
			client = restresourcehandler.NewRestResourceHandler(
				httpClient,
				url,
				restresourcehandler.RestResourceHandlerConfig{
					ResourceEncoding: resourceEncoding,
					RemoteErrorExtractor: func(response *http.Response) error {
						respPayload, err := ioutil.ReadAll(response.Body)

						if err != nil {
							return err
						}

						var remoteError apiError
						err = json.Unmarshal(respPayload, &remoteError)

						if err != nil {
							return err
						}

						return fmt.Errorf(`http status %d, message "%s"`, response.StatusCode, remoteError.ErrorMessage)
					},
				})
		})

		forEachExampleValidApiCall(func(reqName string, req apiCall) {
			It(fmt.Sprintf(`reports custom remote error during "%s" call`, reqName), func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, apiError{"some api error occurred"})))

				err := req(client)

				Expect(err).To(MatchError(fmt.Errorf(`http status 500, message "some api error occurred"`)))
			})
		})
	})
})
