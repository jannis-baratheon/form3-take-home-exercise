package form3apiclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestForm3apiclientModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "form3apiclient testsuite")
}
