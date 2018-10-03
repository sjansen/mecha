package commands

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/sjansen/mecha/internal/config"
	"github.com/sjansen/mecha/internal/fs"
)

type pinCmd struct {
	version string
	add     bool
	remove  bool
}

func (cmd *pinCmd) register(app *kingpin.Application, version string) {
	cmd.version = version

	pin := app.Command("pin", "Create a minimal mecha config or reinitialize an existing one").
		Action(cmd.run)
	pin.Flag("add", "require a specific mecha version").
		Short('a').
		BoolVar(&cmd.add)
	pin.Flag("remove", "stop requiring a specific mecha version").
		Short('r').
		BoolVar(&cmd.remove)
}

func (cmd *pinCmd) run(pc *kingpin.ParseContext) error {
	cfg, err := fs.OpenProjectConfig()
	if err != nil {
		return err
	}

	c := config.Files{Project: cfg}
	if cmd.add || cmd.remove {
		var version string
		if cmd.add {
			version = cmd.version
		}
		before, after := c.SetPinned(version)
		fmt.Println("before:", before)
		fmt.Println("after:", after)
		return cfg.Save()
	}

	version := c.GetPinned()
	if version == "" {
		fmt.Println("status:", "not pinned")
	} else {
		fmt.Println("status:", version)
	}
	return nil
}
