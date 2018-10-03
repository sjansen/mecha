package fs

import (
	"os"
	"path/filepath"

	"github.com/knq/ini"
)

func FindProjectRoot(dir string) (root string, err error) {
	for {
		filename := filepath.Join(dir, ".mecha", "config")
		if _, err := os.Stat(filename); err == nil {
			return dir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if dir == parent {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

func OpenProjectConfig() (f *ini.File, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	root, err := FindProjectRoot(wd)
	if err != nil {
		return
	}

	filename := filepath.Join(root, ".mecha", "config")
	return ini.LoadFile(filename)
}
