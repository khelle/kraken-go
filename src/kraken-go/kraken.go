package main

import (
	"os"
	"./process/wrapper"
	"./errors"
)

func main() {
	process := wrapper.New()
	err := process.Start(os.Args[1:])
	errors.Log(err)

	os.Exit(0)
}
