package main

import (
	"fmt"
	"os"

	"github.com/foks-proj/go-snowpack-compiler/lib"
)

func debugStop() {
	flag := os.Getenv("SNOWPC_DEBUG_STOP")
	if flag == "" || flag == "0" {
		return
	}
	pid := os.Getpid()
	fmt.Fprintf(os.Stderr, "SNOWPC_DEBUG_STOP: pid %d\n", pid)
	fmt.Fprintf(os.Stderr, "Attach debugger and press enter to continue...")
	var buf [1]byte
	// ignore errors / output. since we're just waiting for input in debug
	_, _ = os.Stdin.Read(buf[:])
}

func mainWithErr() error {
	debugStop()

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
