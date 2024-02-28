// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package loggerlevel

import (
	"encoding"
	"fmt"
)

type LoggerLevel string

const (
	DEBUG    LoggerLevel = "DEBUG"
	INFO     LoggerLevel = "INFO"
	WARN     LoggerLevel = "WARN"
	ERROR    LoggerLevel = "ERROR"
	FATAL    LoggerLevel = "FATAL"
	PANIC    LoggerLevel = "PANIC"
	Disabled LoggerLevel = "Disabled"
	TRACE    LoggerLevel = "TRACE"
)

// String returns the string representation of LoggerLevel
func (rcv LoggerLevel) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(LoggerLevel)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for LoggerLevel.
func (rcv *LoggerLevel) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "DEBUG":
		*rcv = DEBUG
	case "INFO":
		*rcv = INFO
	case "WARN":
		*rcv = WARN
	case "ERROR":
		*rcv = ERROR
	case "FATAL":
		*rcv = FATAL
	case "PANIC":
		*rcv = PANIC
	case "Disabled":
		*rcv = Disabled
	case "TRACE":
		*rcv = TRACE
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LoggerLevel`, str)
	}
	return nil
}
