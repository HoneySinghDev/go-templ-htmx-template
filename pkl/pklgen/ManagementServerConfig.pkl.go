// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

import "github.com/apple/pkl-go/pkl"

type ManagementServerConfig struct {
	Secret string `pkl:"secret"`

	ReadinessTimeout *pkl.Duration `pkl:"readinessTimeout"`

	LivenessTimeout *pkl.Duration `pkl:"livenessTimeout"`

	ProbeWritablePathAbs []string `pkl:"probeWritablePathAbs"`

	ProbeWriteableTouchfile string `pkl:"probeWriteableTouchfile"`
}
