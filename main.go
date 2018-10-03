package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/commands"
)

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
