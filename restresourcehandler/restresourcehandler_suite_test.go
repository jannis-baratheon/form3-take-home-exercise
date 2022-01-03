package restresourcehandler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRestresourcehandlerModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "restresourcehandler testsuite")
}
