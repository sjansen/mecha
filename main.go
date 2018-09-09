package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/commands"
)

func main() {
	app := kingpin.
		New("mecha", "A tool to make software development easier").
		UsageTemplate(kingpin.CompactUsageTemplate)

	(&versionCmd{}).register(app)
	commands.Register(app)

	if len(os.Args) == 1 {
		app.Usage(os.Args[1:])
		os.Exit(1)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

type versionCmd struct{}

func (v *versionCmd) register(app *kingpin.Application) {
	app.Command("version", "Print mecha's version").Action(v.run)
}

func (v *versionCmd) run(pc *kingpin.ParseContext) error {
	fmt.Println("mecha", version)
	return nil
}
