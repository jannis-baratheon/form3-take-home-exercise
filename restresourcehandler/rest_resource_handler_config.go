package restresourcehandler

import (
	"fmt"
	"net/http"
	"net/url"
)

type RemoteErrorExtractor func(response *http.Response) error

// TODO maxresponsesize, errordeserializer
type restResourceHandlerConfig struct {
	RemoteErrorExtractor RemoteErrorExtractor
	ResourceURL          url.URL
	ResourceEncoding     string
	DataPropertyName     string
	IsDataWrapped        bool
}

type RestResourceHandlerConfigBuilder interface {
	SetResourceURL(url url.URL) RestResourceHandlerConfigBuilder
	SetResourceEncoding(encoding string) RestResourceHandlerConfigBuilder
	SetDataPropertyName(name string) RestResourceHandlerConfigBuilder
	SetRemoteErrorExtractor(extractor RemoteErrorExtractor) RestResourceHandlerConfigBuilder
	Build() restResourceHandlerConfig
}

type restResourceHandlerConfigBuilder struct {
	RemoteErrorExtractor RemoteErrorExtractor
	ResourceURL          *url.URL
	ResourceEncoding     *string
	DataPropertyName     *string
}

func NewConfigBuilder() RestResourceHandlerConfigBuilder {
	return &restResourceHandlerConfigBuilder{}
}

func (b *restResourceHandlerConfigBuilder) SetResourceURL(url url.URL) RestResourceHandlerConfigBuilder {
	b.ResourceURL = &url
	return b
}

func (b *restResourceHandlerConfigBuilder) SetResourceEncoding(encoding string) RestResourceHandlerConfigBuilder {
	b.ResourceEncoding = &encoding
	return b
}

func (b *restResourceHandlerConfigBuilder) SetDataPropertyName(name string) RestResourceHandlerConfigBuilder {
	b.DataPropertyName = &name
	return b
}

func (b *restResourceHandlerConfigBuilder) SetRemoteErrorExtractor(extractor RemoteErrorExtractor) RestResourceHandlerConfigBuilder {
	b.RemoteErrorExtractor = extractor
	return b
}

func (b *restResourceHandlerConfigBuilder) Build() restResourceHandlerConfig {
	requireNonNil("ResourceURL", b.ResourceURL)
	requireNonNil("ResourceEncoding", b.ResourceEncoding)

	var dataPropertyName string
	if b.DataPropertyName == nil {
		dataPropertyName = ""
	} else {
		dataPropertyName = *b.DataPropertyName
	}

	extractor := b.RemoteErrorExtractor 
	if extractor == nil {
		extractor = func(response *http.Response) error {
			return fmt.Errorf(`remote server returned error status: %d"`, response.StatusCode)
		}
	}

	return restResourceHandlerConfig{
		RemoteErrorExtractor: extractor,
		ResourceURL:          *b.ResourceURL,
		ResourceEncoding:     *b.ResourceEncoding,
		DataPropertyName:     dataPropertyName,
		IsDataWrapped:        b.DataPropertyName != nil,
	}
}

func requireNonNil(property string, v interface{}) {
	if v == nil {
		panic(fmt.Sprintf(`Value of "%s" must not be nil`, property))
	}
}
