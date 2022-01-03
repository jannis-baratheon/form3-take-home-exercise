package form3apiclient

import(
	"testing"
	"fmt"
)

func TestAccountsGet(t *testing.T) {
	client := CreateClient("http://localhost:8080")
	// ad27e265-9605-4b4b-a0e5-3003ea9cc4de
	accounts := client.GetAccounts()
	account := accounts.GetAccount("ad27e265-9605-4b4b-a0e5-3003ea9cc4de")
	fmt.Println("%V", account)
}