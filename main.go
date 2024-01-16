package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"

	"github.com/sjansen/mecha/internal/commands"
)

var build string // set by goreleaser

func main() {
	app := kingpin.
		New("mecha", "A tool to make software development easier").
		UsageTemplate(kingpin.CompactUsageTemplate)
	if build != "" {
		commands.Register(app, build)
	} else {
		commands.Register(app, version)
	}

	if len(os.Args) == 1 {
		app.Usage(os.Args[1:])
		os.Exit(1)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
