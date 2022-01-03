package httpclientrestdecorator_test

import (
	"fmt"
	"github.com/jannis-baratheon/Form3-take-home-excercise/httpclientrestdecorator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/url"
)

type person struct {
	Name string `json:"name"`
}

var _ = Describe("GenericRestClient", func() {
	var server *ghttp.Server
	var client httpclientrestdecorator.GenericRestClient

	BeforeEach(func() {
		server = ghttp.NewServer()
		url, _ := url.Parse(fmt.Sprintf("%s/api", server.URL()))
		client = httpclientrestdecorator.CreateGenericRestClient(httpclientrestdecorator.GenericRestClientConfig{BaseApiUrl: *url})
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
