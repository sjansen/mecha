package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Frob
type FrobCmd struct {
	Filename string
}

func (f *FrobCmd) register(app *kingpin.Application) {
	cmd := app.Command("frob", "Frob a file.").Action(f.run)
	cmd.Arg("FILE", "A filename.").Required().ExistingFileVar(&f.Filename)
}

func (cmd *FrobCmd) run(pc *kingpin.ParseContext) error {
	fmt.Printf("frobbing %q...\n", cmd.Filename)
	return nil
}

// Munge
type MungeCmd struct {
	Filename string
}

func (m *MungeCmd) register(app *kingpin.Application) {
	cmd := app.Command("munge", "Munge a file.").Action(m.run)
	cmd.Arg("FILE", "A filename.").Required().ExistingFileVar(&m.Filename)
}

func (cmd *MungeCmd) run(pc *kingpin.ParseContext) error {
	fmt.Printf("munging %q...\n", cmd.Filename)
	return nil
}

// main
func main() {
	app := kingpin.
		New("kingping-demo", "A tool to explore the kingping API").
		UsageTemplate(kingpin.CompactUsageTemplate)

	f := &FrobCmd{}
	f.register(app)

	m := &MungeCmd{}
	m.register(app)

	if len(os.Args) == 1 {
		app.Usage(os.Args[1:])
		os.Exit(1)
	}

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
