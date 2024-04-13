// Code generated from Pkl module `appConfig.pkl`. DO NOT EDIT.
package pklgen

import (
	"context"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkl/pklgen/environment"
	"github.com/apple/pkl-go/pkl"
)

type Pkl struct {
	Echo *EchoConfig `pkl:"Echo"`

	SB *SupabaseConfig `pkl:"SB"`

	DB *DatabaseConfig `pkl:"DB"`

	Auth *AuthServerConfig `pkl:"Auth"`

	Logger *LoggerConfig `pkl:"Logger"`

	ENV environment.Environment `pkl:"ENV"`

	Management *ManagementServerConfig `pkl:"Management"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Pkl
func LoadFromPath(ctx context.Context, path string) (ret *Pkl, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Pkl
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*Pkl, error) {
	var ret Pkl
	if err := evaluator.EvaluateModule(ctx, source, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}
