package form3apiclient_test

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type wrapper struct {
	AccountData form3apiclient.AccountData `json:"data"`
}

type remoteError struct {
	Message string `json:"error_message"`
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

type apiCall func(client form3apiclient.Form3ApiClient) (interface{}, error)

var exampleValidApiCalls = map[string]apiCall{
	"accounts get": func(client form3apiclient.Form3ApiClient) (interface{}, error) {
		return client.Accounts().Get(uuid.NewString())
	},
	"accounts delete": func(client form3apiclient.Form3ApiClient) (interface{}, error) {
		return client.Accounts().Get(uuid.NewString())
	},
	"accounts create": func(client form3apiclient.Form3ApiClient) (interface{}, error) {
		return client.Accounts().Create(form3apiclient.AccountData{})
	},
}

func forEachExampleValidApiCall(consumer func(string, apiCall)) {
	for reqName, req := range exampleValidApiCalls {
		consumer(reqName, req)
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

	Context("when on happy-path", func() {
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

			response, err := client.Accounts().Get(expectedData.ID)

			Expect(err).To(Succeed())
			Expect(response).To(Equal(expectedData))
		})

		It("should delete account", func() {
			accountId := uuid.NewString()

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", accountsUrl+"/"+accountId, "version=100"),
					ghttp.RespondWith(http.StatusNoContent, nil)))

			err := client.Accounts().Delete(accountId, 100)

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

			actualResponse, err := client.Accounts().Create(requestData)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualResponse).To(Equal(expectedData))
		})

		Context("when remote error occurs", func() {
			forEachExampleValidApiCall(func(callName string, call apiCall) {
				expectedErrorStatus := http.StatusBadRequest
				expectedRemoteErrorMessage := "i have no idea what language you speak"

				It(fmt.Sprintf(`should include server message in error if available for "%s" call`, callName), func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWithJSONEncoded(expectedErrorStatus, remoteError{expectedRemoteErrorMessage})))

					_, err := call(client)

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(
						fmt.Errorf(
							`api responded with error: http status code %d, http status "%d %s", server message: "%s"`,
							expectedErrorStatus,
							expectedErrorStatus,
							http.StatusText(expectedErrorStatus),
							expectedRemoteErrorMessage)))
				})
			})

			forEachExampleValidApiCall(func(callName string, call apiCall) {
				expectedErrorStatus := http.StatusBadRequest

				It(fmt.Sprintf(`should not include server message in error if unavailable for "%s" call`, callName), func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.RespondWith(expectedErrorStatus, nil)))

					_, err := call(client)

					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(
						fmt.Errorf(
							`api responded with error: http status code %d, http status "%d %s"`,
							expectedErrorStatus,
							expectedErrorStatus,
							http.StatusText(expectedErrorStatus))))
				})
			})
		})
	})
})