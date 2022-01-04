package form3apiclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRestresourcehandlerModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "form3apiclient testsuite")
}
