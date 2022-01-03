package restclient_test

import (
	"fmt"
	"github.com/jannis-baratheon/Form3-take-home-excercise/restclient"
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

var _ = Describe("GenericRestClient", func() {
	var server *ghttp.Server
	var httpClient *http.Client
	var client restclient.RestClient

	BeforeEach(func() {
		httpClient = &http.Client{}
		server = ghttp.NewServer()
		url, _ := url.Parse(fmt.Sprintf("%s/api/people", server.URL()))
		config := restclient.
			NewRestClientConfigBuilder().
			SetDataPropertyName("data").
			SetResourceEncoding("application/json").
			SetResourceURL(*url).
			Build()
		client = restclient.CreateRestClient(httpClient, config)
	})

	It("should fetch resource", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/people/1", "attrs=name"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWith(http.StatusOK, `{ "data": {"name": "Smith"} }`)))

		var response person
		err := client.Fetch("1", map[string]string{"attrs": "name"}, &response)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Equal(person{"Smith"}))
	})

	It("should delete resource", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/people/1"),
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
				ghttp.VerifyMimeType("application/json"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyJSONRepresenting(wrapper{payload}),
				ghttp.RespondWithJSONEncoded(http.StatusCreated, wrapper{expectedResponse})))

		var actualResponse person
		err := client.Create(payload, &actualResponse)

		Expect(err).NotTo(HaveOccurred())
		Expect(actualResponse).To(Equal(expectedResponse))
	})

	AfterEach(func() {
		server.Close()
	})
})
