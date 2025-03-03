package main

import (
	"fmt"
	"os"

	"github.com/foks-proj/go-snowpack-compiler/lib"
)

func mainWithErr() error {
	opts, err := lib.ParseOpts()
	if err != nil {
		return err
	}
	err = opts.Run()
	if err != nil {
		return err
	}
	return nil

}

func main() {
	err := mainWithErr()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
