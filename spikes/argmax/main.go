package main

import (
	"fmt"
	"os"

	"github.com/tklauser/go-sysconf"
)

func main() {
	argmax, err := sysconf.Sysconf(sysconf.SC_ARG_MAX)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("argmax =", argmax)

	env := os.Environ()
	argmax -= int64(len(env))
	for _, kv := range env {
		argmax -= int64(len(kv))
	}

	fmt.Println("ajusted argmax =", argmax)
}
