package generator

import (
	"fmt"
	"io"
	"strings"

	"github.com/ettle/strcase"
)

func RenameDestFilename(srcFilename string) string {
	if substr, isCut := strings.CutSuffix(srcFilename, "_test.go"); isCut {
		return substr + "_goacc_test.go"
	} else if substr, isCut := strings.CutSuffix(srcFilename, ".go"); isCut {
		return substr + "_goacc.go"
	}

	panic("filename does not have '.go' ext")
}

func renameToCamelCase(ident string) string {
	return strcase.ToCamel(ident)
}

func renameToPascalCase(ident string) string {
	return strcase.ToPascal(ident)
}

func renameToShorter(ident string) string {

	ident = renameToPascalCase(ident)
	result := make([]rune, 0, len(ident))
	for _, r := range ident {
		if r >= 'A' && r <= 'Z' {
			result = append(result, r)
		}
	}

	return strings.ToLower(string(result))
}

func convertToComment(expr string) string {
	expr, _ = strings.CutSuffix(expr, "\n")
	return "// " + strings.ReplaceAll(expr, "\n", "\n// ")
}

func fprintfln(dest io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(dest, format+"\n", a...); err != nil {
		panic(fmt.Sprintf("failed to write: %s", err))
	}
}
