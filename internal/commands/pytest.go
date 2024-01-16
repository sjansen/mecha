package commands

import (
	"context"
	"time"

	"github.com/alecthomas/kingpin/v2"

	"github.com/sjansen/mecha/internal/pytest"
)

type pytestCmd struct {
	args    []string
	timeout int
}

func (cmd *pytestCmd) register(app *kingpin.Application) {
	pytest := app.Command("pytest", "Run pytest while capturing output and metrics").
		Action(cmd.run)
	pytest.Flag("timeout", "maximum run time in seconds").
		Short('t').Default("60").
		IntVar(&cmd.timeout)
	pytest.Arg("ARGS", "pytest arguments").Required().
		StringsVar(&cmd.args)
}

func (cmd *pytestCmd) run(pc *kingpin.ParseContext) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cmd.timeout)*time.Second,
	)
	defer cancel()

	return pytest.Run(ctx, cmd.args...)
}
