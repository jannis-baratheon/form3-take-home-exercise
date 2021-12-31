package form3apiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValid(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"data\":{\"someprop\": \"someval\"}}"))
		}),
	)
	defer s.Close()

	url, _ := url.Parse(fmt.Sprintf("%s/somerestapi", s.URL))
	client := createGenericRestClient(genericRestClientConfig{baseApiUrl: *url})

	var respJson struct {
		Prop string `json:"someprop"`
	}

	err := client.Fetch("some_resource", "some_id", &respJson)

	fmt.Printf("%v %v", err, respJson)
}
