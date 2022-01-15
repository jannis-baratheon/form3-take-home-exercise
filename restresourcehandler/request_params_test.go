package restresourcehandler

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func someValidRequestParams() requestParams {
	return requestParams{
		HTTPMethod:          "GET",
		DoDiscardResourceID: true,
		DoDiscardContent:    true,
	}
}

var _ = Describe("requestParams", func() {
	Context("panics during validation", func() {
		It("when invalid http method has not been set", func() {
			params := someValidRequestParams()
			params.HTTPMethod = "UNKNOWN_METHOD"

			Expect(func() { validateRequestParameters(params) }).To(PanicWith(`Unknown HTTP method "UNKNOWN_METHOD".`))
		})

		It("when resource id has not been set and it shall not be discarded", func() {
			params := someValidRequestParams()
			params.DoDiscardResourceID = false
			params.ResourceID = ""

			Expect(func() { validateRequestParameters(params) }).To(PanicWith("Invalid request parameters: ResourceID is empty, but DoDiscardResourceID is not set."))
		})

		It("when response content placeholder has not been set and it shall not be discarded", func() {
			params := someValidRequestParams()
			params.DoDiscardContent = false
			params.Response = nil

			Expect(func() { validateRequestParameters(params) }).To(PanicWith("Invalid request parameters: Response is null, but DoDiscardContent is not set."))
		})
	})
})
