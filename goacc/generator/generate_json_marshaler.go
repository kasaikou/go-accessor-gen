package generator

import (
	"io"

	"github.com/kasaikou/goacc/goacc/entity"
)

func generateMarshalJSON(dest io.Writer, structConfig entity.StructConfig) {

	if !structConfig.EnableMarshalJson() {
		return
	}

	recieverName := "__" + renameToShorter(structConfig.Name())
	jsonContentType := renameToPascalCase(structConfig.Name()) + "JSONContent"

	fprintfln(dest, "func (%s *%s) MarshalJSON() ([]byte, error) {", recieverName, structConfig.Name())
	{
		// Defining JSON content struct.
		fprintfln(dest, "")
		fprintfln(dest, "type %s struct {", jsonContentType)
		for _, field := range structConfig.Fields() {
			if field.JsonTag() == "-" {
				continue
			}
			fprintfln(dest, "%s %s `json:\"%s\"`", renameToPascalCase(field.Name()), field.TypeName(), field.JsonTag())
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")

		// Initialize for JSON content struct.
		fprintfln(dest, "return json.Marshal(%s{", jsonContentType)
		for _, field := range structConfig.Fields() {
			if field.JsonTag() == "-" {
				continue
			}
			fprintfln(dest, "%s: %s.%s,", renameToPascalCase(field.Name()), recieverName, field.Name())
		}
		fprintfln(dest, "})")
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")
}
