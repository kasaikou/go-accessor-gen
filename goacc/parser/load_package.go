package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/kasaikou/goacc/goacc/entity"
	"golang.org/x/tools/go/packages"
)

var regGoaccTag = regexp.MustCompile(`(^|,)goacc:"(.*?)"($|,)`)
var regJsonTag = regexp.MustCompile(`(^|,)json\((.*?)\)($|,)`)

func LoadPackage(input *LoadPackageInput) (*packages.Package, error) {

	const mode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
	dir, err := filepath.Abs(input.dirname)
	if err != nil {
		return nil, err
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:  mode,
		Tests: false,
	}, dir)

	if err != nil {
		return nil, err
	} else if len(pkgs) != 1 {
		return nil, fmt.Errorf("%d packages found", len(pkgs))
	} else if len(pkgs[0].GoFiles) != len(pkgs[0].Syntax) {
		return nil, fmt.Errorf("%d compiled go files found, but %d syntax found", len(pkgs[0].CompiledGoFiles), len(pkgs[0].Syntax))
	}

	return pkgs[0], nil
}

func ParsePackage(input *ParsePackageInput) []entity.FileConfig {

	files := []entity.FileConfig{}
	for i := range input.pkg.GoFiles {
		goFile := input.pkg.GoFiles[i]
		syntax := input.pkg.Syntax[i]

		imports := []entity.ImportConfig{}
		structs := []entity.StructConfig{}

		ast.Inspect(syntax, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.ImportSpec:
				name := ""
				if node.Name != nil {
					name = node.Name.Name
				}
				imports = append(imports, *entity.NewImportConfigBuilder(name, node.Path.Value).Purge())
				return false

			case *ast.TypeSpec:
				namedType, structType, isValid := parseNamedStructType(input.pkg.TypesInfo.ObjectOf(node.Name))
				if !isValid {
					break
				}
				structConfig := parseStruct(namedType, structType, input.defaultTag)

				// Set struct metadata into structConfig.
				structConfig.SetDefineFilename(goFile)
				if doc := node.Doc.Text(); doc != "" {
					structConfig.SetDocText(doc)
				} else if comment := node.Comment.Text(); comment != "" {
					structConfig.SetDocText(comment)
				}

				// Set field metadata into FieldConfig.
				structNode := node.Type.(*ast.StructType)
				for i, field := range structNode.Fields.List {
					if doc := field.Doc.Text(); doc != "" {
						structConfig.Fields()[i].SetDocText(doc)
					} else if comment := field.Comment.Text(); comment != "" {
						structConfig.Fields()[i].SetDocText(comment)
					}
					structConfig.Fields()[i].SetTypeName(parseTypeString(field.Type))
				}

				structs = append(structs, structConfig)
				return false
			}
			return true
		})

		files = append(files, *entity.NewFileConfigBuilder(goFile, input.pkg.Name, imports, structs).Purge())
	}

	return files
}

func parseNamedStructType(object types.Object) (namedType *types.Named, structType *types.Struct, isValid bool) {
	namedType, ok := object.Type().(*types.Named)
	if !ok {
		return namedType, structType, false
	}

	structType, ok = namedType.Underlying().(*types.Struct)
	if !ok {
		return namedType, structType, false
	}

	return namedType, structType, true
}

func parseStruct(namedType *types.Named, structType *types.Struct, defaultTag string) entity.StructConfig {

	structSupportsBuilder := entity.NewStructSupportsBuilder()
	enableMarshalJSON := false

	// Check init() and build() support.
	for i := range namedType.NumMethods() {
		method := namedType.Method(i)
		switch method.Name() {
		case "goaccPreNewHook":
			signature := method.Signature()
			if signature.Params().Len() == 0 && signature.Results().Len() == 0 {
				structSupportsBuilder.WithHasPreNewHook(true)
			}
		case "goaccPostNewHook":
			signature := method.Signature()
			if signature.Params().Len() == 0 && signature.Results().Len() == 0 {
				structSupportsBuilder.WithHasPostNewHook(true)
			}
		}
	}

	// Enum member and configurates.
	fields := []entity.FieldConfig{}
	for i := range structType.NumFields() {
		field := structType.Field(i)
		fieldName := field.Name()
		fieldType := field.Type()
		fieldTag := defaultTag
		jsonTag := "-"

		// Parse field tags.
		if parsedTagResult := regGoaccTag.FindStringSubmatch(structType.Tag(i)); len(parsedTagResult) == 4 {
			fieldTag = parsedTagResult[2]
		}

		splitedFieldTags := strings.Split(fieldTag, ",")
		features := parseStructFieldTag(splitedFieldTags)

		// Parse json tag.
		if parsedJsonTagResult := regJsonTag.FindStringSubmatch(fieldTag); len(parsedJsonTagResult) == 4 {
			switch parsedJsonTag := parsedJsonTagResult[2]; parsedJsonTag {
			case "":
				jsonTag = field.Name()
			case ",omitempty":
				jsonTag = field.Name() + ",omitempty"
			default:
				jsonTag = parsedJsonTag
			}
			enableMarshalJSON = true
		} else if slices.Contains(splitedFieldTags, "json") {
			jsonTag = field.Name()
			enableMarshalJSON = true
		}

		fields = append(fields, *entity.NewFieldConfigBuilder(
			fieldName,
			fieldType.String(),
			jsonTag,
			&features,
		).Purge())
	}

	return *entity.NewStructConfigBuilder(
		namedType.Obj().Name(),
		*structSupportsBuilder.Purge(),
		"",
		enableMarshalJSON,
		fields,
	).Purge()
}

func parseStructFieldTag(splitedTags []string) entity.FieldConfigFeatures {
	if len(splitedTags) == 1 && splitedTags[0] == "-" {
		return *entity.NewFieldConfigFeaturesBuilder(false, false, false, false, false, false).Purge()
	}

	return *entity.NewFieldConfigFeaturesBuilder(
		slices.Contains(splitedTags, "mutex"),
		slices.Contains(splitedTags, "required"),
		slices.Contains(splitedTags, "optional"),
		slices.Contains(splitedTags, "getptr"),
		slices.Contains(splitedTags, "get"),
		slices.Contains(splitedTags, "set"),
	).Purge()
}
