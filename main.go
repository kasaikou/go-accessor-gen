package main

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/goacc/goacc/cmd"
	"github.com/kasaikou/goacc/goacc/generator"
)

func main() {
	config := cmd.ParseGenerateFlag(os.Args[1:])
	g := generator.NewGeneratorBuilder(config).Purge()

	wd, _ := os.Getwd()
	doublestar.GlobWalk(os.DirFS(wd), config.IncludePattern(), func(path string, d fs.DirEntry) error {
		if !d.IsDir() {
			if !strings.HasSuffix(path, "_goacc.go") && !strings.HasSuffix(path, "_goacc_test.go") {
				log.Printf("Generate from %s", path)
				generator.WriteFile(g.Generate(path))
			}
		}
		return nil
	})
}
