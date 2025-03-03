package main

import (
	"fmt"
	"os"
)

func mainWithErr() error {
	opts, err := ParseOpts()
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
