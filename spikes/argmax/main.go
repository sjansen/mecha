package main

/*
#include <limits.h>
*/
import "C"

import (
	"fmt"
	"os"
)

func main() {
	argmax := int(uintptr(C.ARG_MAX))
	fmt.Println("argmax =", argmax)

	env := os.Environ()
	argmax -= len(env)
	for _, kv := range env {
		argmax -= len(kv)
	}

	fmt.Println("ajusted argmax =", argmax)
}
