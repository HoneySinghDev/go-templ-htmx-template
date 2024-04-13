// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package environment

import (
	"encoding"
	"fmt"
)

type Environment string

const (
	Dev  Environment = "dev"
	Qa   Environment = "qa"
	Prod Environment = "prod"
)

// String returns the string representation of Environment
func (rcv Environment) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(Environment)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Environment.
func (rcv *Environment) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "dev":
		*rcv = Dev
	case "qa":
		*rcv = Qa
	case "prod":
		*rcv = Prod
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid Environment`, str)
	}
	return nil
}
