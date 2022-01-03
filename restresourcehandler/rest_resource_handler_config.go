package restresourcehandler

import (
	"net/url"
)

// TODO maxresponsesize, errordeserializer
type RestResourceHandlerConfig struct {
	RemoteErrorExtractor RemoteErrorExtractor
	ResourceURL          url.URL
	ResourceEncoding     string
	DataPropertyName     string
	IsDataWrapped        bool
}

func validateRestResourceHandlerConfig(config RestResourceHandlerConfig) {
	if config.IsDataWrapped && config.DataPropertyName == "" {
		panic("IsDataWrapped is set, but DataPropertyName has not been given.")
	}

	if !config.IsDataWrapped && config.DataPropertyName != "" {
		panic("IsDataWrapped is not set, but DataPropertyName has been given.")
	}

	if !config.ResourceURL.IsAbs() {
		panic("Resource URL must be absolute.")
	}

	if config.ResourceEncoding == "" {
		panic("ResourceEncoding must be set.")
	}
}
