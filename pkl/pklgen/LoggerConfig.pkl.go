// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

import "github.com/HoneySinghDev/go-templ-htmx-template/pkl/pklgen/loggerlevel"

type LoggerConfig struct {
	Level loggerlevel.LoggerLevel `pkl:"level"`

	RequestLevel loggerlevel.LoggerLevel `pkl:"requestLevel"`

	RequestBody bool `pkl:"requestBody"`

	RequestHeader bool `pkl:"requestHeader"`

	RequestQuery bool `pkl:"requestQuery"`

	ResponseHeader bool `pkl:"responseHeader"`

	ResponseBody bool `pkl:"responseBody"`

	LogCaller bool `pkl:"logCaller"`

	PreetyPrintConsole bool `pkl:"preetyPrintConsole"`
}
