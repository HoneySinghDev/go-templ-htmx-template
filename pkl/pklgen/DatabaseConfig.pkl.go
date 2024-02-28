// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

import "github.com/apple/pkl-go/pkl"

type DatabaseConfig struct {
	PSQLDB string `pkl:"PSQLDB"`

	PSQLHOST string `pkl:"PSQLHOST"`

	PSQLPORT int32 `pkl:"PSQLPORT"`

	PSQLUSER string `pkl:"PSQLUSER"`

	PSQLPASS string `pkl:"PSQLPASS"`

	AdditionalParams map[string]string `pkl:"AdditionalParams"`

	DBMaxOpenConns int `pkl:"DBMaxOpenConns"`

	MaxIdleConns int `pkl:"MaxIdleConns"`

	MinIdleConns int `pkl:"MinIdleConns"`

	ConnectionMaxLifetime *pkl.Duration `pkl:"ConnectionMaxLifetime"`
}
