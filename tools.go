//go:build tools

package tools

// Tooling dependencies
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md

// This file may incorporate tools that may be *both* used as CLI and as lib
// Keep in mind that these global tools change the go.mod/go.sum dependency tree
// Other tooling may be installed as *static binary* directly within the Dockerfile

import (
	_ "github.com/cosmtrek/air"
	_ "github.com/jdudmesh/gomon"
	_ "github.com/pressly/goose/v3"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
