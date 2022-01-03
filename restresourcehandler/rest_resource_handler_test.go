package restresourcehandler_test

import (
	"fmt"
	"github.com/jannis-baratheon/Form3-take-home-excercise/restresourcehandler"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/url"
)

type person struct {
	Name string `json:"name"`
}

type wrapper struct {
	Data person `json:"data"`
}

type request func(client restresourcehandler.RestResourceHandler) error

var exampleValidRequests = map[string]request {
	"fetch": func (client restresourcehandler.RestResourceHandler) error {
		var response person
		return client.Fetch("1", map[string]string{"attrs": "name"}, &response)
	},
	"delete": func (client restresourcehandler.RestResourceHandler) error {
		return client.Delete("1", map[string]string{"version": "1"})
	},
	"create": func (client restresourcehandler.RestResourceHandler) error {
		var actualResponse person
		return client.Create(person{"Smith"}, &actualResponse)
	},
}

var _ = Describe("RestResourceHandler", func() {
	var server *ghttp.Server
	var httpClient *http.Client
	var client restresourcehandler.RestResourceHandler
	const resourceEncoding = "application/json"

	BeforeEach(func() {
		httpClient = &http.Client{}
		server = ghttp.NewServer()
		url, _ := url.Parse(fmt.Sprintf("%s/api/people", server.URL()))
		config := restresourcehandler.RestResourceHandlerConfig{
			IsDataWrapped:    true,
			DataPropertyName: "data",
			ResourceEncoding: resourceEncoding,
			ResourceURL:      *url,
		}
		client = restresourcehandler.NewRestResourceHandler(httpClient, config)
	})

	AfterEach(func() {
		server.Close()
	})

	It("should fetch resource", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/people/1", "attrs=name"),
				ghttp.VerifyHeaderKV("Accept", resourceEncoding),
				ghttp.RespondWith(http.StatusOK, `{ "data": {"name": "Smith"} }`)))

		var response person
		err := client.Fetch("1", map[string]string{"attrs": "name"}, &response)

		Expect(err).To(Succeed())
		Expect(response).To(Equal(person{"Smith"}))
	})

	It("should delete resource", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/people/1", "version=1"),
				ghttp.RespondWith(http.StatusNoContent, nil)))

		err := client.Delete("1", map[string]string{"version": "1"})

		Expect(err).NotTo(HaveOccurred())
	})

	It("should create resource", func() {
		payload := person{"Smith"}
		expectedResponse := person{"Gennings"}

		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/people"),
				ghttp.VerifyMimeType(resourceEncoding),
				ghttp.VerifyHeaderKV("Accept", resourceEncoding),
				ghttp.VerifyJSONRepresenting(wrapper{payload}),
				ghttp.RespondWithJSONEncoded(http.StatusCreated, wrapper{expectedResponse})))

		var actualResponse person
		err := client.Create(payload, &actualResponse)

		Expect(err).NotTo(HaveOccurred())
		Expect(actualResponse).To(Equal(expectedResponse))
	})

	for reqName, req := range exampleValidRequests {
		It(fmt.Sprintf("should report remote error during %s", reqName), func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusInternalServerError, nil)))

			err := req(client)

			Expect(err).To(MatchError(fmt.Errorf("remote server returned error status: 500")))
		})
	}
})
