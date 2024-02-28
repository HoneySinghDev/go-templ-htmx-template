// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

type EchoConfig struct {
	Debug bool `pkl:"debug"`

	ListenAddr string `pkl:"listenAddr"`

	HideInternalServerErrorDetails bool `pkl:"hideInternalServerErrorDetails"`

	BaseUrl *string `pkl:"baseUrl"`

	LoggerMiddleware bool `pkl:"loggerMiddleware"`

	RecoverMiddleware bool `pkl:"recoverMiddleware"`

	SecureMiddleware *EchoServerSecureMiddlewareConfig `pkl:"secureMiddleware"`
}
