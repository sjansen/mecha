package fs

import (
	"os"
	"path/filepath"

	"github.com/knq/ini"
)

func OpenConfig() (f *ini.File, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		filename := filepath.Join(dir, ".mecha", "config")
		if _, err := os.Stat(filename); err == nil {
			return ini.LoadFile(filename)
		}

		parent := filepath.Dir(dir)
		if dir == parent {
			break
		}

		dir = parent
	}
	return nil, nil
}
