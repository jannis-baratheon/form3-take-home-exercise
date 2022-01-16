package form3apiclient

import (
	"net/http"
)

type form3APIRemoteError struct {
	ErrorMessage string `json:"error_message"`
}

// Form3ApiClient is a client object used to call the Form3 REST API.
type Form3ApiClient struct {
	accountsEndpoint *accounts
}

// NewForm3APIClient constructs a Form3 API Client for the given URL (e.g. "http://localhost:8080/v1")
// and HTTP client instance.
//
// All HTTP calls will be made using the passed in HTTP client.
func NewForm3APIClient(apiURL string, httpClient *http.Client) *Form3ApiClient {
	accounts, err := newAccounts(apiURL, httpClient)
	if err != nil {
		panic(err)
	}

	return &Form3ApiClient{accountsEndpoint: accounts}
}

// Accounts returns a handler for the accounts endpoint  of the Form3 REST API
// ("< form3 api url>/organisation/accounts").
func (c *Form3ApiClient) Accounts() *accounts {
	return c.accountsEndpoint
}
