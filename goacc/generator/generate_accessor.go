package generator

import (
	"io"
	"strings"

	"github.com/kasaikou/goacc/goacc/entity"
)

func generateAccessor(dest io.Writer, structConfig entity.StructConfig) {

	ptrGetterFields := make([]entity.FieldConfig, 0, len(structConfig.Fields()))
	getterFields := make([]entity.FieldConfig, 0, len(structConfig.Fields()))
	setterFields := make([]entity.FieldConfig, 0, len(structConfig.Fields()))

	for _, field := range structConfig.Fields() {
		if field.Features().HasPtrGetter() {
			ptrGetterFields = append(ptrGetterFields, field)
		}

		if field.Features().HasGetter() && field.Name()[0] >= 'a' && field.Name()[0] <= 'z' {
			getterFields = append(getterFields, field)
		}

		if field.Features().HasSetter() {
			setterFields = append(setterFields, field)
		}
	}

	structType := structConfig.Name()
	recieverName := "__" + renameToShorter(structType)

	for _, field := range getterFields {
		getterName := renameToPascalCase(field.Name())
		if comment := field.DocText(); comment != "" {
			fprintfln(dest, strings.ReplaceAll(convertToComment(comment), field.Name(), getterName))
		}
		fprintfln(dest, "func (%s *%s) %s() %s {", recieverName, structType, getterName, field.TypeName())
		{
			fprintfln(dest, "if %s != nil {", recieverName)
			{
				// TODO: add mutex
				fprintfln(dest, "return %s.%s", recieverName, field.Name())
			}
			fprintfln(dest, "}")
			fprintfln(dest, "")

			fprintfln(dest, "panic(\"%s is nil\")", structType)
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")
	}

	for _, field := range ptrGetterFields {
		getterName := renameToPascalCase(field.Name()) + "Ptr"
		if comment := field.DocText(); comment != "" {
			fprintfln(dest, strings.ReplaceAll(convertToComment(comment), field.Name(), getterName))
		}
		fprintfln(dest, "func (%s *%s) %s() *%s {", recieverName, structType, getterName, field.TypeName())
		{
			fprintfln(dest, "if %s != nil {", recieverName)
			{
				// TODO: add mutex
				fprintfln(dest, "return &%s.%s", recieverName, field.Name())
			}
			fprintfln(dest, "}")
			fprintfln(dest, "")

			fprintfln(dest, "panic(\"%s is nil\")", structType)
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")
	}

	for _, field := range setterFields {
		fprintfln(dest, "func (%s *%s) Set%s(%s %s) {", recieverName, structType, renameToPascalCase(field.Name()), field.Name(), field.TypeName())
		{
			fprintfln(dest, "if %s != nil {", recieverName)
			{
				// TODO: add mutex
				fprintfln(dest, "%s.%s = %s", recieverName, field.Name(), field.Name())
			}
			fprintfln(dest, "}")
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")
	}
}
