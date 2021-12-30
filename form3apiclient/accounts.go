package form3apiclient

import(
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"encoding/json"
)

type AccountsResource interface {
	GetAccount(id string) AccountData
}

type accountsResource struct {
	Context context
}

func createAccountResource(context context) AccountsResource {
	return &accountsResource{Context: context}
}

func (accountsResource accountsResource) GetAccount(id string) AccountData {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/organisation/accounts/%s", accountsResource.Context.Url, id), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var respJson struct {
		Data AccountData `json:"data"`
	}
	err = json.Unmarshal(body, &respJson)

	if err != nil { 
		log.Fatal("Failed to parse response. ", err)
	}

	return respJson.Data
}
