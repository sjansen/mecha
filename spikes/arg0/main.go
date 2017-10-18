package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var arg0 struct {
	err error
	val string
}

func init() {
	arg0 := os.Args[0]
	if filepath.IsAbs(arg0) {
		setARG0(arg0, nil)
	} else if strings.ContainsRune(arg0, os.PathSeparator) {
		setARG0(filepath.Abs(arg0))
	} else {
		setARG0(exec.LookPath(arg0))
	}
}

func setARG0(path string, err error) {
	if err == nil {
		path, err = filepath.EvalSymlinks(path)
	}
	arg0.err = err
	arg0.val = path
}

func main() {
	fmt.Printf("%q (%v)\n", arg0.val, arg0.err)
}
