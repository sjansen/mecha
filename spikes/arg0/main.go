package main

import (
	"fmt"
	"os"
)

func main() {
	path, err := os.Executable()
	fmt.Printf("%q (%v)\n", path, err)
}
