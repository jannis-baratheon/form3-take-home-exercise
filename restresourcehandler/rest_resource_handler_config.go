package restresourcehandler

import (
	"fmt"
	"net/url"
)

// TODO maxresponsesize, errordeserializer
type restResourceHandlerConfig struct {
	ResourceURL      url.URL
	ResourceEncoding string
	DataPropertyName string
	IsDataWrapped    bool
}

type RestResourceHandlerConfigBuilder interface {
	SetResourceURL(url url.URL) RestResourceHandlerConfigBuilder
	SetResourceEncoding(encoding string) RestResourceHandlerConfigBuilder
	SetDataPropertyName(name string) RestResourceHandlerConfigBuilder
	Build() restResourceHandlerConfig
}

type restResourceHandlerConfigBuilder struct {
	ResourceURL      *url.URL
	ResourceEncoding *string
	DataPropertyName *string
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

func (b *restResourceHandlerConfigBuilder) Build() restResourceHandlerConfig {
	requireNonNil("ResourceURL", b.ResourceURL)
	requireNonNil("ResourceEncoding", b.ResourceEncoding)

	var dataPropertyName string
	if b.DataPropertyName == nil {
		dataPropertyName = ""
	} else {
		dataPropertyName = *b.DataPropertyName
	}

	return restResourceHandlerConfig{
		ResourceURL:      *b.ResourceURL,
		ResourceEncoding: *b.ResourceEncoding,
		DataPropertyName: dataPropertyName,
		IsDataWrapped:    b.DataPropertyName != nil,
	}
}

func requireNonNil(property string, v interface{}) {
	if v == nil {
		panic(fmt.Sprintf(`Value of "%s" must not be nil`, property))
	}
}
