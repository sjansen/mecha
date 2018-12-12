package commands

import (
	"math/rand"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/tui"
)

type startCmd struct{}

func (cmd *startCmd) register(app *kingpin.Application) {
	app.Command("start", "Start the application defined by Procfile").
		Action(cmd.run)
}

func (cmd *startCmd) run(pc *kingpin.ParseContext) error {
	updates := make(chan *tui.Status)
	screen := tui.NewStackedTextViews()
	screen.AddStatusItem("todo", "TODO:", updates)

	go func() {
		for {
			if ok := rand.Intn(100) > 20; ok {
				updates <- &tui.Status{
					Severity: tui.Healthy,
					Message:  "PASS",
				}
			} else {
				updates <- &tui.Status{
					Severity: tui.Alert,
					Message:  "FAIL",
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return screen.Run()
}
