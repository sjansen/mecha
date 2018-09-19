package commands

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type versionCmd struct {
	version string
}

func (v *versionCmd) register(app *kingpin.Application, version string) {
	v.version = version
	app.Command("version", "Print mecha's version").Action(v.run)
}

func (v *versionCmd) run(pc *kingpin.ParseContext) error {
	fmt.Println("mecha", v.version)
	return nil
}
