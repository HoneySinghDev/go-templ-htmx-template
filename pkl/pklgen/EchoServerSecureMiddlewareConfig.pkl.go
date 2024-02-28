// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

type EchoServerSecureMiddlewareConfig struct {
	Enable bool `pkl:"enable"`

	XssProtection string `pkl:"xssProtection"`

	ContentTypeNosniff string `pkl:"contentTypeNosniff"`

	XFrameOptions string `pkl:"xFrameOptions"`

	HstsMaxAge int `pkl:"hstsMaxAge"`

	HstsExcludeSubdomains bool `pkl:"hstsExcludeSubdomains"`

	HstsPreload bool `pkl:"hstsPreload"`

	ContentSecurityPolicy string `pkl:"contentSecurityPolicy"`

	CspReportOnly bool `pkl:"cspReportOnly"`

	ReferrerPolicy string `pkl:"referrerPolicy"`
}
