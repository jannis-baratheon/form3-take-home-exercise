package restresourcehandler

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func someValidRequestParams() requestParams {
	return requestParams{
		HttpMethod:          "GET",
		DoDiscardResourceId: true,
		DoDiscardContent:    true,
	}
}

var _ = FDescribe("requestParams", func() {
	Context("panics during validation", func() {
		It("when invalid http method has not been set", func() {
			params := someValidRequestParams()
			params.HttpMethod = "UNKNOWN_METHOD"

			Expect(func() { validateRequestParameters(params) }).To(PanicWith(`Unknown HTTP method "UNKNOWN_METHOD".`))
		})

		It("when resource id has not been set and it shall not be discarded", func() {
			params := someValidRequestParams()
			params.DoDiscardResourceId = false
			params.ResourceId = ""

			Expect(func() { validateRequestParameters(params) }).To(PanicWith("Invalid request parameters: ResourceId is empty, but DoDiscardResourceId is not set."))
		})

		It("when response content placeholder has not been set and it shall not be discarded", func() {
			params := someValidRequestParams()
			params.DoDiscardContent = false
			params.Response = nil

			Expect(func() { validateRequestParameters(params) }).To(PanicWith("Invalid request parameters: Response is null, but DoDiscardContent is not set."))
		})
	})
})
