package form3apiclient_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
	"net/http"
	"testing"
)

func TestAPICLient(t *testing.T) {
	accounts := form3apiclient.NewForm3APIClient("http://localhost:8080/v1", &http.Client{}).Accounts()

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
	var response form3apiclient.AccountData
	response, err := accounts.Create(accountData)

	fmt.Printf("\nCREATE ***************\n\nerror: %v\n\nresponse: %v\n\n", err, response)

	if err != nil {
		t.Fail()
	}

	account, err := accounts.Get(accountData.ID)

	fmt.Printf("\nFETCH ****************\n\nerror: %v\n\nresponse: %v\n\n", err, account)

	if err != nil {
		t.Fail()
	}

	err = accounts.Delete(account.ID, 0)

	fmt.Printf("\nDELETE ***************\n\nerror: %v\n\n", err)

	if err != nil {
		t.Fail()
	}
}
