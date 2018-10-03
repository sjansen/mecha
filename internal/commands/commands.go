package commands

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func Register(app *kingpin.Application, version string) {
	(&initCmd{}).register(app)
	(&pinCmd{}).register(app, version)
	(&versionCmd{}).register(app, version)
}
