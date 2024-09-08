package cmd

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kasaikou/goacc/goacc/entity"
	"github.com/kasaikou/goacc/goacc/generator"
)

func Generate(config *entity.GenerateConfig) {
	g := generator.NewGenerator()
	doublestar.GlobWalk(os.DirFS(config.WorkingDir()), config.IncludePattern(), func(path string, d fs.DirEntry) error {
		if !d.IsDir() {
			if !strings.HasSuffix(path, "_goacc.go") && !strings.HasSuffix(path, "_goacc_test.go") {
				log.Printf("Generate from %s", path)
				generator.WriteFile(g.Generate(path, config))
			}
		}
		return nil
	})
}
