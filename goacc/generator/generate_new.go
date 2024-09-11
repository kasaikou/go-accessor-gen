package generator

import (
	"io"
	"strings"

	"github.com/kasaikou/goacc/goacc/entity"
)

func generateNew(dest io.Writer, structConfig entity.StructConfig) {

	requiredFields := make([]entity.FieldConfig, 0, len(structConfig.Fields()))
	optionalFields := make([]entity.FieldConfig, 0, len(structConfig.Fields()))
	for _, field := range structConfig.Fields() {
		if field.Features().HasRequired() {
			requiredFields = append(requiredFields, field)
		}
		if field.Features().HasOptional() {
			optionalFields = append(optionalFields, field)
		}
	}

	if len(requiredFields) == 0 && len(optionalFields) == 0 {
		return
	}

	structType := structConfig.Name()
	builderType := renameToPascalCase(structType) + "Builder"
	identFieldStruct := "__" + renameToShorter(structType)

	// Creates PiyoPiyoBuilder struct definition.
	fprintfln(dest, convertToComment(`
		%s is an instance for generating an instance of %s.
	`), builderType, structType)
	fprintfln(dest, "type %s struct {\n", builderType)
	{
		fprintfln(dest, "%s *%s", identFieldStruct, structType)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

	// Creates NewPiyoPiyoBuilder() function.
	fprintfln(dest, convertToComment(`
		New%s creates an %s instance.
	`), builderType, builderType)
	fprintfln(dest, "func New%s(", builderType)
	for _, field := range requiredFields {
		fprintfln(dest, "%s %s,", field.Name(), field.TypeName())
	}
	fprintfln(dest, ") *%s {", builderType)
	{
		fprintfln(dest, "%s := &%s{}", identFieldStruct, structType)
		fprintfln(dest, "")
		if structConfig.StructSupportsPtr().HasPreNewHook() {
			fprintfln(dest, "%s.goaccPreNewHook() // This function calls your defined hook.", identFieldStruct)
			fprintfln(dest, "")
		}

		// TODO: add mutex

		if len(requiredFields) > 0 {
			for _, field := range requiredFields {
				fprintfln(dest, "%s.%s = %s", identFieldStruct, field.Name(), field.Name())
			}
			fprintfln(dest, "")
		}

		fprintfln(dest, "return &%s{%s: %s}", builderType, identFieldStruct, identFieldStruct)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

	recieverName := identFieldStruct + "b"

	// Creates (PiyoPiyoBuilder).WithHogeHoge() functions.
	for _, field := range optionalFields {
		withName := "Set" + renameToPascalCase(field.Name())
		if comment := field.DocText(); comment != "" {
			fprintfln(dest, strings.ReplaceAll(convertToComment(comment), field.Name(), withName))
		}
		fprintfln(dest, "func (%s *%s) %s(%s %s) *%s {", recieverName, builderType, withName, field.Name(), field.TypeName(), builderType)
		{
			fprintfln(dest, "if %s == nil {", recieverName)
			{
				fprintfln(dest, "panic(\"%s is nil\")", builderType)
			}
			fprintfln(dest, "} else if %s.%s != nil {", recieverName, identFieldStruct)
			{
				// TODO: add mutex

				fprintfln(dest, "%s.%s.%s = %s", recieverName, identFieldStruct, field.Name(), field.Name())
				fprintfln(dest, "return %s", recieverName)
			}
			fprintfln(dest, "}")
			fprintfln(dest, "")

			fprintfln(dest, "panic(\"%s has been already purged\")", structType)
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")
	}

	// Creates (PiyoPiyoBuilder).Purge() function.
	fprintfln(dest, convertToComment(`
		Purge purges %s instance from %s.

		If calls other method in %s after Purge called, it will be panic.
	`), structType, builderType, builderType)
	fprintfln(dest, "func (%s *%s) Build() *%s {", recieverName, builderType, structType)
	{
		fprintfln(dest, "if %s == nil {", recieverName)
		{
			fprintfln(dest, "panic(\"%s is nil\")", builderType)
		}
		fprintfln(dest, "} else if %s.%s != nil {", recieverName, identFieldStruct)
		{
			fprintfln(dest, "%s := %s.%s", identFieldStruct, recieverName, identFieldStruct)
			fprintfln(dest, "%s.%s = nil", recieverName, identFieldStruct)
			fprintfln(dest, "")

			if structConfig.StructSupportsPtr().HasPostNewHook() {
				fprintfln(dest, "%s.goaccPostNewHook() // This function calls your defined hook.", identFieldStruct)
				fprintfln(dest, "")
			}

			fprintfln(dest, "return %s", identFieldStruct)
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")

		fprintfln(dest, "panic(\"%s has been already purged\")", structType)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

}
