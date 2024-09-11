package cmd

import (
	"flag"
	"os"

	"github.com/kasaikou/goacc/goacc/entity"
)

func ParseGenerateFlag(args []string) *entity.GenerateConfig {
	wd, _ := os.Getwd()
	includePattern := flag.String("i", "**.go", "Search file pattern.")
	defaultTag := flag.String("t", "-", "Default label, uses cannot found 'goacc' label.")
	// workingDirectory := flag.String("d", wd, "Working directory path.")

	flag.CommandLine.Parse(args)

	return entity.NewGenerateConfigBuilder(
		wd,
		*includePattern,
		*defaultTag,
	).Build()
}
