package restresourcehandler

import "net/http"

// RemoteErrorExtractor is a function prototype for functions
// which extract additional data from an error server response.
type RemoteErrorExtractor func(response *http.Response) error

// Config represents configuration of a REST API endpoint.
type Config struct {
	// RemoteErrorExtractor is a function that should extract
	// additional data from an error response.
	RemoteErrorExtractor RemoteErrorExtractor
	// ResourceEncoding denotes the encoding to be used when
	// encoding (request) or decoding (response) server data.
	// Currently only JSON encodings are supported.
	ResourceEncoding string
	// IsDataWrapped denotes if the DTOs in server responses
	// should be deserialized from the root of the JSON response (false)
	// or are rather nested in the response JSON
	// (true, e.g. a HATEOAS response wraps the DTO in a root "data" property).
	IsDataWrapped bool
	// DataPropertyName is the property name in the response JSON
	// in which the response DTO should be looked for
	// (in case IsDataWrapped is true).
	DataPropertyName string
}

// validateRestResourceHandlerConfig does a sanity check of a Config instance.
func validateRestResourceHandlerConfig(config Config) {
	if config.IsDataWrapped && config.DataPropertyName == "" {
		panic("IsDataWrapped is set, but DataPropertyName has not been given.")
	}

	if !config.IsDataWrapped && config.DataPropertyName != "" {
		panic("IsDataWrapped is not set, but DataPropertyName has been given.")
	}

	if config.ResourceEncoding == "" {
		panic("ResourceEncoding must be set.")
	}
}
