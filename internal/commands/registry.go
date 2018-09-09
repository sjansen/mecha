package commands

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func Register(app *kingpin.Application) {
	(&initCmd{}).register(app)
}
