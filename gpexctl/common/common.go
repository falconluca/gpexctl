package common

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

type GlobalFlags struct {
	UITable bool
	Debug   bool
}

var (
	Flags GlobalFlags
)

func ExitWithErrorf(format string, a ...interface{}) {
	ExitWithError(fmt.Errorf(format, a...))
}

func ExitWithError(err error) {
	if err != nil {
		color.New(color.FgRed).Fprint(os.Stderr, "Error: ")
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
