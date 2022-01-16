package form3apiclient

import (
	"fmt"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
)

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type AccountData struct {
	Attributes     AccountAttributes `json:"attributes,omitempty"`
	ID             string            `json:"id,omitempty"`
	OrganisationID string            `json:"organisation_id,omitempty"`
	Type           string            `json:"type,omitempty"`
	Version        int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   string   `json:"account_classification,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 string   `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            bool     `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  string   `json:"status,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
}

type Accounts interface {
	Get(id string) (AccountData, error)
	Delete(id string, version int64) error
	Create(accountData AccountData) (AccountData, error)
}

type accounts struct {
	Handler restresourcehandler.RestResourceHandler
}

const resourcePath = "organisation/accounts"

func newAccounts(apiURL string, httpClient *http.Client) (*accounts, error) {
	accountsResourceURL, err := join(apiURL, resourcePath)
	if err != nil {
		return nil, WrapError(err, "constructing api url")
	}

	handler := restresourcehandler.NewRestResourceHandler(httpClient, accountsResourceURL, config)

	return &accounts{handler}, nil
}

func (a *accounts) Get(id string) (AccountData, error) {
	var accountData AccountData
	err := a.Handler.Fetch(id, nil, &accountData)

	return accountData, err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}

func (a *accounts) Delete(id string, version int64) error {
	err := a.Handler.Delete(id, map[string]string{"version": fmt.Sprint(version)})

	return err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}

func (a *accounts) Create(accountData AccountData) (AccountData, error) {
	var response AccountData
	err := a.Handler.Create(&accountData, &response)

	return response, err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}
