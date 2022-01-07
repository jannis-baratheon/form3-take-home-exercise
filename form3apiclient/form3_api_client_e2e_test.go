package form3apiclient_test

import (
	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Form3ApiClient with real server", Label("e2e"), Ordered, func() {
	var accounts form3apiclient.Accounts

	BeforeEach(func() {
		url := os.Getenv("FORM3_API_URL")

		if url == "" {
			panic("FORM3_API_URL has to be set")
		}
		accounts = form3apiclient.NewForm3APIClient(url, &http.Client{}).Accounts()
	})

	It("runs simple pipeline", func() {
		var resource form3apiclient.AccountData

		By("create account", func() {
			accountData := form3apiclient.AccountData{
				ID:             uuid.NewString(),
				OrganisationID: uuid.NewString(),
				Type:           "accounts",
				Attributes: form3apiclient.AccountAttributes{
					AccountClassification: "Personal",
					Name:                  []string{"Jan Kowalski"},
					Country:               "PL",
				},
			}
			var err error
			resource, err = accounts.Create(accountData)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				By("delete account", func() {
					err := accounts.Delete(resource.ID, 0)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		By("fetch account", func() {
			fetchedAccountData, err := accounts.Get(resource.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedAccountData).To(Equal(resource))
		})
	})
})
