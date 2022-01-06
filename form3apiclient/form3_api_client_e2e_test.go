package form3apiclient_test

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form3ApiClient with real server", Label("e2e"), func() {
	var accounts form3apiclient.Accounts
	url := os.Getenv("FORM3_API_URL")

	if url == "" {
		panic("FORM3_API_URL has to be set")
	}

	BeforeEach(func() {
		accounts = form3apiclient.NewForm3APIClient(url, &http.Client{}).Accounts()
	})

	It("runs simple pipeline without error", func() {
		var resource form3apiclient.AccountData

		By("creating an account", func() {
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
				By("and then deleting it", func() {
					err := accounts.Delete(resource.ID, 0)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		By("fetching it afterwards", func() {
			fetchedAccountData, err := accounts.Get(resource.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedAccountData).To(Equal(resource))
		})
	})
})
