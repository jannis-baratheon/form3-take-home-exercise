package restclient

import (
	"fmt"
	"net/url"
)

// TODO maxresponsesize, errordeserializer
type restClientConfig struct {
	ResourceURL      url.URL
	ResourceEncoding string
	DataPropertyName string
	IsDataWrapped    bool
}

type RestClientConfigBuilder interface {
	SetResourceURL(url url.URL) RestClientConfigBuilder
	SetResourceEncoding(encoding string) RestClientConfigBuilder
	SetDataPropertyName(name string) RestClientConfigBuilder
	Build() restClientConfig
}

type restClientConfigBuilder struct {
	ResourceURL      *url.URL
	ResourceEncoding *string
	DataPropertyName *string
}

func NewRestClientConfigBuilder() RestClientConfigBuilder {
	return &restClientConfigBuilder{}
}

func (b *restClientConfigBuilder) SetResourceURL(url url.URL) RestClientConfigBuilder {
	b.ResourceURL = &url
	return b
}

func (b *restClientConfigBuilder) SetResourceEncoding(encoding string) RestClientConfigBuilder {
	b.ResourceEncoding = &encoding
	return b
}

func (b *restClientConfigBuilder) SetDataPropertyName(name string) RestClientConfigBuilder {
	b.DataPropertyName = &name
	return b
}

func (b *restClientConfigBuilder) Build() restClientConfig {
	requireNonNil("ResourceURL", b.ResourceURL)
	requireNonNil("ResourceEncoding", b.ResourceEncoding)

	var dataPropertyName string
	if b.DataPropertyName == nil {
		dataPropertyName = ""
	} else {
		dataPropertyName = *b.DataPropertyName
	}

	return restClientConfig{
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
