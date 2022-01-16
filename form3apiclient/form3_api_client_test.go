package form3apiclient_test

import (
	"fmt"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	someValidUUID      = "ad27e265-9605-4b4b-a0e5-3003ea9cc422"
	someOtherValidUUID = "ad27e265-9605-4b4b-a0e5-3003ea9cc422"
	accountsURL        = "/organisation/accounts"
)

type wrapper struct {
	AccountData form3apiclient.AccountData `json:"data"`
}

type remoteError struct {
	Message string `json:"error_message"`
}

type apiCall func(client form3apiclient.Form3ApiClient) error

var exampleValidAPICalls = map[string]apiCall{
	"accounts get": func(client form3apiclient.Form3ApiClient) error {
		_, err := client.Accounts().Get(someValidUUID)

		return err //nolint:wrapcheck // we need this error unwrapped
	},
	"accounts delete": func(client form3apiclient.Form3ApiClient) error {
		return client.Accounts().Delete(someValidUUID, 0) //nolint:wrapcheck // we need this error unwrapped
	},
	"accounts create": func(client form3apiclient.Form3ApiClient) error {
		_, err := client.Accounts().Create(form3apiclient.AccountData{})

		return err //nolint:wrapcheck // we need this error unwrapped
	},
}

func forEachExampleValidAPICall(consumer func(string, apiCall)) {
	for reqName, req := range exampleValidAPICalls {
		consumer(reqName, req)
	}
}

var _ = Describe("Form3ApiClient", func() {
	var server *ghttp.Server
	var client form3apiclient.Form3ApiClient

	const resourceEncoding = "application/json; charset=utf-8"

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = form3apiclient.NewForm3APIClient(server.URL(), &http.Client{})
	})

	AfterEach(func() {
		server.Close()
	})

	Context("on happy-path", func() {
		It("gets account", func() {
			expectedData := someValidAccountData(someValidUUID)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", accountsURL+"/"+expectedData.ID),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.RespondWithJSONEncoded(http.StatusOK, wrapper{expectedData})))

			response, err := client.Accounts().Get(expectedData.ID)

			Expect(err).To(Succeed())
			Expect(response).To(Equal(expectedData))
		})

		It("deletes account", func() {
			accountID := someValidUUID

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", accountsURL+"/"+accountID, "version=100"),
					ghttp.RespondWith(http.StatusNoContent, nil)))

			err := client.Accounts().Delete(accountID, 100)

			Expect(err).NotTo(HaveOccurred())
		})

		It("creates account", func() {
			requestData := someValidAccountData(someValidUUID)
			expectedData := someValidAccountData(someOtherValidUUID)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", accountsURL),
					ghttp.VerifyContentType(resourceEncoding),
					ghttp.VerifyHeaderKV("Accept", resourceEncoding),
					ghttp.VerifyJSONRepresenting(wrapper{requestData}),
					ghttp.RespondWithJSONEncoded(http.StatusCreated, wrapper{expectedData})))

			actualResponse, err := client.Accounts().Create(requestData)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualResponse).To(Equal(expectedData))
		})
	})

	Context("when remote error occurs", func() {
		Context("and server provides an error message", func() {
			expectedErrorStatus := http.StatusBadRequest
			expectedRemoteErrorMessage := "i have no idea what language you speak"

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(expectedErrorStatus, remoteError{expectedRemoteErrorMessage})))
			})

			forEachExampleValidAPICall(func(callName string, call apiCall) {
				It(fmt.Sprintf(`includes server message in returned error for "%s" call`, callName), func() {
					err := call(client)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(expectedRemoteErrorMessage))
				})
			})
		})

		Context("and server does not provide an error message", func() {
			expectedErrorStatus := http.StatusBadRequest

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(expectedErrorStatus, nil)))
			})

			forEachExampleValidAPICall(func(callName string, call apiCall) {
				It(fmt.Sprintf(`does not include the server message part in returned error for "%s" call`, callName), func() {
					err := call(client)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).NotTo(ContainSubstring("server message: "))
				})
			})
		})
	})
})
