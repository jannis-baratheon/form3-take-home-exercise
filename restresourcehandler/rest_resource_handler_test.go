package restresourcehandler_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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

var exampleValidAPICalls = map[string]apiCall{
	"fetch": func(client restresourcehandler.RestResourceHandler) error {
		var response person

		return client.Fetch("1", map[string]string{"attrs": "name"}, &response) //nolint:wrapcheck,lll // we need this error unwrapped
	},
	"delete": func(client restresourcehandler.RestResourceHandler) error {
		return client.Delete("1", map[string]string{"version": "1"}) //nolint:wrapcheck // we need this error unwrapped
	},
	"create": func(client restresourcehandler.RestResourceHandler) error {
		var response person

		return client.Create(person{"Smith"}, &response) //nolint:wrapcheck // we need this error unwrapped
	},
}

func forEachExampleValidAPICall(consumer func(string, apiCall)) {
	for callName, call := range exampleValidAPICalls {
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

	Context("on happy-path", func() {
		var client restresourcehandler.RestResourceHandler

		BeforeEach(func() {
			client = restresourcehandler.NewRestResourceHandler(
				httpClient,
				url,
				restresourcehandler.Config{
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
				restresourcehandler.Config{
					ResourceEncoding: resourceEncoding,
				})
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil)))
		})

		forEachExampleValidAPICall(func(reqName string, req apiCall) {
			It(fmt.Sprintf(`provides default error during "%s" call`, reqName), func() {
				err := req(client)

				Expect(err).To(MatchError(restresourcehandler.RemoteError(http.StatusInternalServerError)))
			})
		})
	})

	Context("with custom remote error extractor", func() {
		Context("providing error not based on response content", func() {
			var client restresourcehandler.RestResourceHandler
			customError := fmt.Errorf("some custom error") //nolint:goerr113 // not a problem here

			BeforeEach(func() {
				client = restresourcehandler.NewRestResourceHandler(
					httpClient,
					url,
					restresourcehandler.Config{
						ResourceEncoding: resourceEncoding,
						RemoteErrorExtractor: func(response *http.Response) error {
							return customError
						},
					})
			})

			forEachExampleValidAPICall(func(reqName string, req apiCall) {
				It(fmt.Sprintf(`reports custom remote error during "%s" call`, reqName), func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWith(http.StatusInternalServerError, nil)))

					err := req(client)

					Expect(err).To(MatchError(customError))
				})
			})
		})

		Context("providing error based on response content", func() {
			var client restresourcehandler.RestResourceHandler

			BeforeEach(func() {
				client = restresourcehandler.NewRestResourceHandler(
					httpClient,
					url,
					restresourcehandler.Config{
						ResourceEncoding: resourceEncoding,
						RemoteErrorExtractor: func(response *http.Response) error {
							respPayload, err := ioutil.ReadAll(response.Body)
							if err != nil {
								panic(err)
							}

							var remoteError apiError
							err = json.Unmarshal(respPayload, &remoteError)

							if err != nil {
								panic(err)
							}

							return restresourcehandler.RemoteErrorWithServerMessage(response.StatusCode, remoteError.ErrorMessage) //nolint:wrapcheck,lll
						},
					})
			})

			forEachExampleValidAPICall(func(reqName string, req apiCall) {
				It(fmt.Sprintf(`reports custom remote error during "%s" call`, reqName), func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, apiError{"some api error occurred"})))

					err := req(client)

					Expect(err).To(MatchError(restresourcehandler.RemoteErrorWithServerMessage(500, "some api error occurred")))
				})
			})
		})
	})
})
