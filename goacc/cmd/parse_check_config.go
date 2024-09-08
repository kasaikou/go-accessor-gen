package cmd

import (
	"flag"
	"os"

	"github.com/kasaikou/goacc/goacc/entity"
)

func ParseCheckFlag(args []string) *entity.CheckConfig {
	wd, _ := os.Getwd()
	includePattern := flag.String("i", "**.go", "Search file pattern.")

	flag.CommandLine.Parse(args)

	return entity.NewCheckConfigBuilder(
		wd,
		*includePattern,
	).Purge()
}
