package cmd

import (
	"flag"

	"github.com/kasaikou/goacc/goacc/entity"
)

func ParseGenerateFlag(args []string) *entity.GenerateConfig {
	includePattern := flag.String("i", "**.go", "Search file pattern.")
	defaultTag := flag.String("t", "-", "Default label, uses cannot found 'goacc' label.")

	flag.CommandLine.Parse(args)

	return entity.NewGenerateConfigBuilder(
		*includePattern,
		*defaultTag,
	).Purge()
}
