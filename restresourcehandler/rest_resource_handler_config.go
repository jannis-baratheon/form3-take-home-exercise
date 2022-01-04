package restresourcehandler

type RestResourceHandlerConfig struct {
	RemoteErrorExtractor RemoteErrorExtractor
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

	if config.ResourceEncoding == "" {
		panic("ResourceEncoding must be set.")
	}
}
