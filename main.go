package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"strings"

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
			if !strings.HasSuffix(path, "_goacc.go") && !strings.HasSuffix(path, "_goacc_test.go") {
				log.Printf("Generate from %s", path)
				generator.WriteFile(g.Generate(path))
			}
		}
		return nil
	})
}
