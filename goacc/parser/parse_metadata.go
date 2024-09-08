package parser

import (
	"bufio"
	"io"
	"strings"

	"github.com/kasaikou/goacc/goacc/entity"
)

func ParseMetadata(r io.Reader) *entity.Metadata {
	builder := entity.NewMetadataBuilder()
	const (
		prefixDefaultTag = "// defaultTag="
		prefixPackage    = "package "
	)

	for s := bufio.NewScanner(r); s.Scan(); {
		line := s.Text()
		if strings.HasPrefix(line, prefixPackage) {
			return builder.Purge()

		} else if defaultTag, isCut := strings.CutPrefix(line, prefixDefaultTag); isCut {
			builder.WithDefaultTag(defaultTag)

		}
	}

	panic("reached EOF before package definition")
}
