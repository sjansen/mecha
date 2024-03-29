package commands

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

type initCmd struct{}

const template = `[core]
config_version = 0
`

func (cmd *initCmd) register(app *kingpin.Application) {
	app.Command(
		"init", "Create a minimal mecha config or reinitialize an existing one",
	).Action(cmd.run)
}

func (cmd *initCmd) run(pc *kingpin.ParseContext) (err error) {
	if err = os.MkdirAll(".mecha", os.ModePerm); err != nil {
		return
	}

	filename := filepath.Join(".mecha", "config")
	if err = os.WriteFile(filename, []byte(template), os.ModePerm); err != nil {
		return
	}

	return nil
}
