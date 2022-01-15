package form3apiclient_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form3ApiClient with real server", Label("e2e"), Ordered, func() {
	var accounts form3apiclient.Accounts
	var createdResources []form3apiclient.AccountData
	var apiUrl string

	createAndScheduleCleanup := func(accountData form3apiclient.AccountData) (form3apiclient.AccountData, error) {
		res, err := accounts.Create(accountData)

		if err == nil {
			createdResources = append(createdResources, res)
		}

		return res, err
	}

	cleanup := func() {
		httpClient := &http.Client{}
		resourceBaseUrl, err := url.Parse(apiUrl)
		if err != nil {
			panic("api url is not a valid url")
		}
		resourceBaseUrl.Path = path.Join(resourceBaseUrl.Path, "/organisation/accounts")

		for _, resource := range createdResources {
			deleteRequest, err := http.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("%s/%s?version=%d", resourceBaseUrl.String(), resource.ID, resource.Version),
				nil)
			if err != nil {
				log.Println(fmt.Sprintf(`WARNING: Failed to cleanup AccountData resource after test. Error: "%v"`, err))
				continue
			}

			resp, err := httpClient.Do(deleteRequest)

			if err != nil || (resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound) {
				log.Println(fmt.Sprintf(`WARNING: Failed to cleanup AccountData resource after test. HTTP Status: "%v", error: "%v"`, resp.Status, err))
			}
		}
	}

	BeforeEach(func() {
		apiUrl = os.Getenv("FORM3_API_URL")

		if apiUrl == "" {
			panic("FORM3_API_URL has to be set")
		}
		accounts = form3apiclient.NewForm3APIClient(apiUrl, &http.Client{}).Accounts()

		DeferCleanup(cleanup)
	})

	It("executes basic happy-path scenario", func() {
		var resource form3apiclient.AccountData

		By("creates account", func() {
			accountData := form3apiclient.AccountData{
				ID:             uuid.NewString(),
				OrganisationID: uuid.NewString(),
				Type:           "accounts",
				Attributes: form3apiclient.AccountAttributes{
					AccountClassification: "Personal",
					Name:                  []string{"Jan Kowalski", "Jasiu Kowalski"},
					Country:               "PL",
				},
			}
			var err error
			resource, err = createAndScheduleCleanup(accountData)
			Expect(err).NotTo(HaveOccurred())
		})

		By("fetches account", func() {
			fetchedAccountData, err := accounts.Get(resource.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedAccountData).To(Equal(resource))
			Expect(fetchedAccountData).To(HaveField("Type", "accounts"))
			Expect(fetchedAccountData).To(HaveField("Attributes.AccountClassification", "Personal"))
			Expect(fetchedAccountData).To(HaveField("Attributes.Country", "PL"))
			Expect(fetchedAccountData).To(HaveField("Attributes.Name", ConsistOf("Jan Kowalski", "Jasiu Kowalski")))
		})

		By("deletes account", func() {
			err := accounts.Delete(resource.ID, resource.Version)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("communicates api errors", func() {
		It("when deleting non-existent account", func() {
			var err error
			var accountData form3apiclient.AccountData

			By("making sure the account does not exist first", func() {
				accountData, err = createAndScheduleCleanup(someValidAccountData(uuid.NewString()))
				Expect(err).NotTo(HaveOccurred())

				err = accounts.Delete(accountData.ID, accountData.Version)
				Expect(err).NotTo(HaveOccurred())
			})

			err = accounts.Delete(accountData.ID, accountData.Version)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf(`api responded with error: http status code 404, http status "404 Not Found"`)))
		})

		It("when attempting to delete account with invalid version", func() {
			var err error
			var accountData form3apiclient.AccountData

			By("making sure the account does exist", func() {
				accountData, err = createAndScheduleCleanup(someValidAccountData(uuid.NewString()))
				Expect(err).NotTo(HaveOccurred())
			})

			err = accounts.Delete(accountData.ID, accountData.Version+1)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf(`api responded with error: http status code 409, http status "409 Conflict", server message: "invalid version"`)))
		})

		It("when creating account with invalid data", func() {
			invalidAccountData := someValidAccountData(uuid.NewString())
			invalidAccountData.Attributes.Name = nil

			_, err := createAndScheduleCleanup(invalidAccountData)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf("api responded with error: http status code 400, http status \"400 Bad Request\", server message: \"validation failure list:\nvalidation failure list:\nvalidation failure list:\nname in body is required\"")))
		})

		It("creating a duplicate account", func() {
			var err error
			var accountData form3apiclient.AccountData

			By("making sure the account does exist", func() {
				accountData, err = createAndScheduleCleanup(someValidAccountData(uuid.NewString()))
				Expect(err).NotTo(HaveOccurred())
			})

			_, err = createAndScheduleCleanup(accountData)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf(`api responded with error: http status code 409, http status "409 Conflict", server message: "Account cannot be created as it violates a duplicate constraint"`)))
		})

		It("fetching a non-existent account", func() {
			var err error
			var accountData form3apiclient.AccountData

			By("making sure the account does not exist first", func() {
				accountData, err = createAndScheduleCleanup(someValidAccountData(uuid.NewString()))
				Expect(err).NotTo(HaveOccurred())

				err = accounts.Delete(accountData.ID, accountData.Version)
				Expect(err).NotTo(HaveOccurred())
			})

			_, err = accounts.Get(accountData.ID)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fmt.Errorf(`api responded with error: http status code 404, http status "404 Not Found", server message: "record ` + accountData.ID + ` does not exist"`)))
		})
	})
})
