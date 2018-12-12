package commands

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/tui"
)

type startCmd struct{}

func (cmd *startCmd) register(app *kingpin.Application) {
	app.Command("start", "Start the application defined by Procfile").
		Action(cmd.run)
}

func (cmd *startCmd) run(pc *kingpin.ParseContext) error {
	screen := tui.NewStackedTextViews()
	screen.AddStatusItem("todo", "TODO:")

	return screen.Run()
}
