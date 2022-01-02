package form3apiclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenericrestclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generic Rest Client Suite")
}
