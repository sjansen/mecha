package commands

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

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

	before := cfg.GetKey("core.version")
	if !(cmd.add || cmd.remove) {
		if before == "" {
			fmt.Println("status:", "not pinned")
		} else {
			fmt.Println("status:", before)
		}
		return nil
	}

	var after string
	if cmd.add {
		if before == cmd.version {
			after = "no change"
		} else {
			after = cmd.version
			cfg.SetKey("core.version", after)
		}
	} else if cmd.remove {
		if before == "" {
			after = "no change"
		} else {
			after = "not pinned"
			cfg.RemoveKey("core.version")
		}
	}

	if err = cfg.Save(); err != nil {
		return err
	}

	if before == "" {
		before = "not pinned"
	}
	fmt.Println("before:", before)
	fmt.Println("after:", after)

	return nil
}
