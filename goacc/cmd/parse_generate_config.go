package cmd

import (
	"flag"

	"github.com/kasaikou/goacc/goacc/entity"
)

func ParseGenerateFlag(args []string) *entity.GenerateConfig {
	includePattern := flag.String("i", "**.go", "Search file pattern.")
	defaultLabel := flag.String("l", "-", "Default label, uses cannot found 'goacc' label.")

	flag.CommandLine.Parse(args)

	return entity.NewGenerateConfigBuilder(
		*includePattern,
		*defaultLabel,
	).Purge()
}
