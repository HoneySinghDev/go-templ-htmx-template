// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

import "github.com/apple/pkl-go/pkl"

type AuthServerConfig struct {
	AccessTokenValidity *pkl.Duration `pkl:"accessTokenValidity"`

	PasswordresetTokenValidity *pkl.Duration `pkl:"passwordresetTokenValidity"`

	DefaultUserScopes []string `pkl:"defaultUserScopes"`

	LastAuthenticatedAtThreshold *pkl.Duration `pkl:"lastAuthenticatedAtThreshold"`
}
