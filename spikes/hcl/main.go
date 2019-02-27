package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl"
)

type Config struct {
	Project string
	Workers int
	Tasks   []*Task `hcl:"task,expand"`
}

type Task struct {
	Name string `hcl:",key"`
	Args []string
}

var configFile = `
project = "foo"
workers = 42

task "bar" {
    args = ["date", "-u"]
}

task "baz" {
    args = ["echo", "Spoon!"]
}
`

func main() {
	var cfg *Config
	err := hcl.Decode(&cfg, configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pprint := spew.ConfigState{
		Indent:                  "    ",
		DisableCapacities:       true,
		DisablePointerAddresses: true,
	}
	pprint.Dump(cfg)
}
