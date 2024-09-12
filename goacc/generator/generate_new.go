package generator

import (
	"fmt"
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
	builderImplType := renameToCamelCase(builderType) + "Impl"
	identFieldStruct := "__" + renameToShorter(structType)
	recieverName := identFieldStruct + "b"

	fprintfln(dest, "type %s interface {", builderType)
	{
		for _, field := range optionalFields {
			_, _, definition := generateNewSetDefinition(builderType, &field)
			fprintfln(dest, definition)
		}
		_, definition := generateNewBuildDefinition(&structConfig)
		fprintfln(dest, definition)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

	// Creates piyoPiyoBuilderImpl struct definition.
	fprintfln(dest, convertToComment(`
		%s is an instance for generating an instance of %s.
	`), builderImplType, structType)
	fprintfln(dest, "type %s struct {\n", builderImplType)
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
	fprintfln(dest, ") %s {", builderType)
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

		fprintfln(dest, "return &%s{%s: %s}", builderImplType, identFieldStruct, identFieldStruct)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

	// Creates (PiyoPiyoBuilder).SetHogeHoge() functions.
	for _, field := range optionalFields {
		methodName, paramName, methodDefinition := generateNewSetDefinition(builderType, &field)
		if comment := field.DocText(); comment != "" {
			fprintfln(dest, strings.ReplaceAll(convertToComment(comment), field.Name(), methodName))
		}
		fprintfln(dest, "func (%s *%s) %s {", recieverName, builderImplType, methodDefinition)
		{
			fprintfln(dest, "if %s == nil {", recieverName)
			{
				fprintfln(dest, "panic(\"%s is nil\")", builderImplType)
			}
			fprintfln(dest, "} else if %s.%s != nil {", recieverName, identFieldStruct)
			{
				// TODO: add mutex

				fprintfln(dest, "%s.%s.%s = %s", recieverName, identFieldStruct, field.Name(), paramName)
				fprintfln(dest, "return %s", recieverName)
			}
			fprintfln(dest, "}")
			fprintfln(dest, "")

			fprintfln(dest, "panic(\"%s has been already purged\")", structType)
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")
	}

	// Creates (PiyoPiyoBuilder).Build() function.
	buildMethodName, buildMethodDefinition := generateNewBuildDefinition(&structConfig)
	fprintfln(dest, convertToComment(`
		%s purges %s instance from %s.

		If calls other method in %s after Purge called, it will be panic.
	`), buildMethodName, structType, builderImplType, builderImplType)

	fprintfln(dest, "func (%s *%s) %s {", recieverName, builderImplType, buildMethodDefinition)

	{
		fprintfln(dest, "if %s == nil {", recieverName)
		{
			fprintfln(dest, "panic(\"%s is nil\")", builderImplType)
		}
		fprintfln(dest, "} else if %s.%s != nil {", recieverName, identFieldStruct)
		{
			fprintfln(dest, "%s := %s.%s", identFieldStruct, recieverName, identFieldStruct)
			fprintfln(dest, "%s.%s = nil", recieverName, identFieldStruct)
			fprintfln(dest, "")

			if structConfig.StructSupportsPtr().HasPostNewHookError() {
				fprintfln(dest, "err := %s.goaccPostNewHook() // This function calls your defined hook.", identFieldStruct)
				fprintfln(dest, "return %s, err", identFieldStruct)

			} else {
				if structConfig.StructSupportsPtr().HasPostNewHook() {
					fprintfln(dest, "%s.goaccPostNewHook() // This function calls your defined hook.", identFieldStruct)
					fprintfln(dest, "")
				}
				fprintfln(dest, "return %s", identFieldStruct)

			}
		}
		fprintfln(dest, "}")
		fprintfln(dest, "")

		fprintfln(dest, "panic(\"%s has been already purged\")", structType)
	}
	fprintfln(dest, "}")
	fprintfln(dest, "")

}

func generateNewSetDefinition(builderType string, config *entity.FieldConfig) (methodName, paramName, definition string) {
	methodName = "Set" + renameToPascalCase(config.Name())
	paramName = config.Name()
	return methodName, paramName, fmt.Sprintf("%s(%s %s) %s", methodName, paramName, config.TypeName(), builderType)
}

func generateNewBuildDefinition(config *entity.StructConfig) (methodName, definition string) {
	if config.StructSupportsPtr().HasPostNewHookError() {
		return "Build", fmt.Sprintf("Build() (*%s, error)", config.Name())
	} else {
		return "Build", fmt.Sprintf("Build() *%s", config.Name())
	}
}
