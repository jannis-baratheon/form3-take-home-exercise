package form3apiclient_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"
)

func TestAPICLient(t *testing.T) {
	client := form3apiclient.NewClient("http://localhost:8080/v1", &http.Client{})

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
	response, err := client.CreateAccount(accountData)

	fmt.Printf("\nCREATE ***************\n\nerror: %v\n\nresponse: %v\n\n", err, response)

	if err != nil {
		t.Fail()
	}

	account, err := client.GetAccount(accountData.ID)

	fmt.Printf("\nFETCH ****************\n\nerror: %v\n\nresponse: %v\n\n", err, account)

	if err != nil {
		t.Fail()
	}

	err = client.DeleteAccount(account.ID, 0)

	fmt.Printf("\nDELETE ***************\n\nerror: %v\n\n", err)

	if err != nil {
		t.Fail()
	}
}
