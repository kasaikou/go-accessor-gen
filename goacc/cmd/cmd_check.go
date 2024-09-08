package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/kasaikou/goacc/goacc/entity"
	"github.com/kasaikou/goacc/goacc/generator"
	"github.com/kasaikou/goacc/goacc/parser"
)

func Check(config *entity.CheckConfig) {
	g := generator.NewGenerator()
	doublestar.GlobWalk(os.DirFS(config.WorkingDir()), config.IncludePattern(), func(path string, d fs.DirEntry) error {
		if !d.IsDir() {
			if !strings.HasSuffix(path, "_goacc.go") && !strings.HasSuffix(path, "_goacc_test.go") {
				goaccPath := generator.RenameDestFilename(path)
				log.Printf("Check with %s", goaccPath)
				f, err := os.Open(goaccPath)

				var current []byte
				var meta *entity.Metadata
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						current = nil
						meta = entity.NewMetadataBuilder().Purge()
					} else {
						panic(err)
					}
				} else {
					buf := bytes.NewBuffer([]byte{})
					defer f.Close()
					meta = parser.ParseMetadata(io.TeeReader(f, buf))
					current = buf.Bytes()
				}

				_, expect := g.Generate(path, entity.NewGenerateConfigBuilder(
					config.WorkingDir(),
					config.IncludePattern(),
					meta.DefaultTag(),
				).Purge())

				switch {
				case expect == nil && current == nil:
				case expect == nil && current != nil:
					fmt.Printf("%s is should not be existed\n", goaccPath)
				case expect != nil && current == nil:
					fmt.Printf("%s is should be generated\n", goaccPath)
				default:
					if diff := cmp.Diff(string(current), string(expect)); diff != "" {
						fmt.Printf("%s is different with expected.\n", goaccPath)
						fmt.Println(diff)
						os.Exit(1)
					}
				}
			}
		}
		return nil
	})
}
