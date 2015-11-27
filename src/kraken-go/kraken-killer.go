package main

import (
	"os"
//	"io"
//	"bytes"
	"fmt"
	"./errors"
)

func main() {

	process, err := os.FindProcess(4432)
	if err != nil {
		errors.Log(errors.New(21, err.Error()))
	}

	fmt.Printf("%#v\n", process.Stdin)

//	stdin, err := process.StdinPipe()
//	if err != nil {
//		errors.Log(errors.New(22, err.Error()))
//	}

//	io.CopyN(stdin, bytes.NewBufferString("woah!\n"), 4096)

//	process.Kill()

	os.Exit(0)
}
