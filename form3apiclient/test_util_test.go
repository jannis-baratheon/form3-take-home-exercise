package form3apiclient_test

import "github.com/jannis-baratheon/Form3-take-home-excercise/form3apiclient"

func someValidAccountData(id string) form3apiclient.AccountData {
	return form3apiclient.AccountData{
		ID:             id,
		OrganisationID: someValidUuid,
		Type:           "accounts",
		Attributes: form3apiclient.AccountAttributes{
			AccountClassification: "Personal",
			Name:                  []string{"Jan Kowalski"},
			Country:               "PL",
		},
	}
}