package commands

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func Register(app *kingpin.Application, version string) {
	(&initCmd{}).register(app)
	(&pinCmd{}).register(app, version)
	(&pytestCmd{}).register(app)
	(&startCmd{}).register(app)
	(&versionCmd{}).register(app, version)
}
