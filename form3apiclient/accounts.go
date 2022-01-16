package form3apiclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jannis-baratheon/form3-take-home-exercise/restresourcehandler"
)

// Accounts allows fetching and modifying account data hosted in a the application.
type Accounts interface {
	// Get fetches account data for the given account id.
	// Context can be used to control asynchronous requests.
	Get(ctx context.Context, id string) (AccountData, error)

	// Delete deletes an account  with the given id and version.
	// Context can be used to control asynchronous requests.
	Delete(ctx context.Context, id string, version int64) error

	// Create creates an account using the passed in AccountData DTO instance.
	// Returns the created account instance.
	// Context can be used to control asynchronous requests.
	Create(ctx context.Context, accountData AccountData) (AccountData, error)
}

func (a *accounts) Get(ctx context.Context, accountID string) (AccountData, error) {
	var accountData AccountData
	err := a.Handler.Fetch(ctx, accountID, nil, &accountData)

	return accountData, err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}

func (a *accounts) Delete(ctx context.Context, accountID string, version int64) error {
	err := a.Handler.Delete(ctx, accountID, map[string]string{"version": fmt.Sprint(version)})

	return err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}

// Create creates an account using the passed in AccountData DTO instance.
// Returns the created account instance.
// Context can be used to control asynchronous requests.
func (a *accounts) Create(ctx context.Context, accountData AccountData) (AccountData, error) {
	var response AccountData
	err := a.Handler.Create(ctx, &accountData, &response)

	return response, err //nolint:wrapcheck // this error is in fact local (see extractRemoteError)
}

type accounts struct {
	Handler *restresourcehandler.RestResourceHandler
}

const resourcePath = "organisation/accounts"

func newAccounts(apiURL string, httpClient *http.Client) (*accounts, error) {
	accountsResourceURL, err := join(apiURL, resourcePath)
	if err != nil {
		return nil, WrapError(err, "constructing api url")
	}

	handler := restresourcehandler.NewRestResourceHandler(
		httpClient, accountsResourceURL, getRestResourceHandlerConfig())

	return &accounts{handler}, nil
}
