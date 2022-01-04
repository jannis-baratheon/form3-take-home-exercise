package form3apiclient_test

import (
	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

type wrapper struct {
	AccountData form3apiclient.AccountData `json:"data"`
}

func minimalValidAccountData() form3apiclient.AccountData {
	return form3apiclient.AccountData{
		ID:             uuid.NewString(),
		OrganisationID: uuid.NewString(),
		Type:           "accounts",
		Attributes: form3apiclient.AccountAttributes{
			AccountClassification: "Personal",
			Name:                  []string{"Jan Kowalski"},
			Country:               "PL",
		},
	}
}

var _ = Describe("Form3ApiClient", func() {
	var server *ghttp.Server
	var httpClient *http.Client
	var apiUrl string
	var accountsUrl string

	const resourceEncoding = "application/json; charset=utf-8"

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiUrl = server.URL()
		accountsUrl = "/organisation/accounts"
		httpClient = &http.Client{}
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when request is valid", func() {
		var client form3apiclient.Form3ApiClient

		BeforeEach(func() {
			client = form3apiclient.NewForm3APIClient(apiUrl, httpClient)
		})

		It("should fetch account", func() {
			expectedData := minimalValidAccountData()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", accountsUrl+"/"+expectedData.ID),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.RespondWithJSONEncoded(http.StatusOK, wrapper{expectedData})))

			response, err := client.GetAccount(expectedData.ID)

			Expect(err).To(Succeed())
			Expect(response).To(Equal(expectedData))
		})

		It("should delete account", func() {
			accountId := uuid.NewString()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", accountsUrl+"/"+accountId, "version=100"),
					ghttp.RespondWith(http.StatusNoContent, nil)))

			err := client.DeleteAccount(accountId, 100)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should create account", func() {
			requestData := minimalValidAccountData()
			expectedData := minimalValidAccountData()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", accountsUrl),
					ghttp.VerifyContentType(resourceEncoding),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.VerifyJSONRepresenting(wrapper{requestData}),
					ghttp.RespondWithJSONEncoded(http.StatusCreated, wrapper{expectedData})))

			actualResponse, err := client.CreateAccount(requestData)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualResponse).To(Equal(expectedData))
		})
	})
})
