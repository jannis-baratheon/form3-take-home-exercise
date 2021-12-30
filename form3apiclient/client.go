package form3apiclient

type Form3ApiClient interface {
	GetAccounts() AccountsResource
}

type form3ApiClient struct {
	Context  context
	Accounts AccountsResource
}

func CreateClient(url string) Form3ApiClient {
	context := createContext(url)
	accountsResource := createAccountResource(context)
	return &form3ApiClient{Context: context, Accounts: accountsResource}
}

func (client form3ApiClient) GetAccounts() AccountsResource {
	return client.Accounts
}
