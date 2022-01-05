package form3apiclient_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestForm3apiclientModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "form3apiclient testsuite")
}
