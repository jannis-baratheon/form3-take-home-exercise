package restresourcehandler

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func someValidRestResourceHandlerConfig() RestResourceHandlerConfig {
	return RestResourceHandlerConfig{
		ResourceEncoding: "some encoding",
	}
}

var _ = Describe("RestResourceHandlerConfig", func() {
	Context("should panic when validated", func() {
		It("when data is wrapped but no property name has been given", func() {
			config := someValidRestResourceHandlerConfig()
			config.IsDataWrapped = true
			config.DataPropertyName = ""

			Expect(func() { validateRestResourceHandlerConfig(config) }).To(PanicWith("IsDataWrapped is set, but DataPropertyName has not been given."))
		})

		It("when data is not wrapped but property name has been given", func() {
			config := someValidRestResourceHandlerConfig()
			config.IsDataWrapped = false
			config.DataPropertyName = "someproperty"

			Expect(func() { validateRestResourceHandlerConfig(config) }).To(PanicWith("IsDataWrapped is not set, but DataPropertyName has been given."))
		})

		It("when resource enoding is not set", func() {
			config := someValidRestResourceHandlerConfig()
			config.ResourceEncoding = ""

			Expect(func() { validateRestResourceHandlerConfig(config) }).To(PanicWith("ResourceEncoding must be set."))
		})
	})
})
