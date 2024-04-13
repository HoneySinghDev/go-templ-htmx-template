package main

import (
	"embed"

	"github.com/HoneySinghDev/go-templ-htmx-template/cmd"
)

//go:embed static
var static embed.FS

func main() {
	cmd.App(static)
}
