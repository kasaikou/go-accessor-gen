package main

import (
	"os"

	"github.com/kasaikou/goacc/goacc/cmd"
)

func main() {
	switch os.Args[1] {
	case "check":
		cmd.Check(cmd.ParseCheckFlag(os.Args[2:]))
	default:
		cmd.Generate(cmd.ParseGenerateFlag(os.Args[1:]))
	}
}
