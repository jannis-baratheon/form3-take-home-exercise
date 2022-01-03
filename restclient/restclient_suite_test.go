package restclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRestclientModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "restclient testsuite")
}
