package main

import (
	"flag"
	"io/fs"
	"os"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/goacc/goacc/generator"
)

func main() {
	g := generator.NewGenerator()

	includePattern := flag.String("i", "**.go", "Search file pattern (default: '**.go')")

	flag.Parse()

	wd, _ := os.Getwd()
	doublestar.GlobWalk(os.DirFS(wd), *includePattern, func(path string, d fs.DirEntry) error {
		if !d.IsDir() {
			generator.WriteFile(g.Generate(path))
		}
		return nil
	})
}
