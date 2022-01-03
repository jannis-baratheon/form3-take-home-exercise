package httpclientrestdecorator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHttpclientrestdecorator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HTTP client REST decorator suite")
}
