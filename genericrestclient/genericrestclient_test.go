package form3apiclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/jannis-baratheon/Form3-take-home-excercise/genericrestclient"
	"net/url"
	"net/http"
	"fmt"
)

type person struct {
	Name string `json:"name"`
}

var _ = Describe("GenericRestClient", func() {
	var server *ghttp.Server
	var client form3apiclient.GenericRestClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		url, _ := url.Parse(fmt.Sprintf("%s/api", server.URL()))
		client = form3apiclient.CreateGenericRestClient(form3apiclient.GenericRestClientConfig{BaseApiUrl: *url})
	})

	It("should fetch resource", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/people/1"),
				ghttp.RespondWith(http.StatusOK, `{ "data": {"name": "Smith"} }`)))

		var response person
		err := client.Fetch("people", "1", &response)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Equal(person{"Smith"}))
	})

	AfterEach(func() {
		server.Close()
	})
})
