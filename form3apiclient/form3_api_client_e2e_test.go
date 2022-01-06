package form3apiclient_test

import (
	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = PDescribe("Form3ApiClient with real server", Label("e2e"), func() {
	var accounts form3apiclient.Accounts

	BeforeEach(func() {
		accounts = form3apiclient.NewForm3APIClient("http://localhost:8080/v1", &http.Client{}).Accounts()
	})

	It("runs simple pipeline without error", func() {
		var accountData form3apiclient.AccountData

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
			_, err := accounts.Create(accountData)
			Expect(err).NotTo(HaveOccurred())
		})

		By("fetching it afterwards", func() {
			fetchedAccountData, err := accounts.Get(accountData.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedAccountData).To(Equal(accountData))
		})

		By("and finally deleting it", func() {
			err := accounts.Delete(accountData.ID, 0)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
