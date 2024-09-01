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

var regTag = regexp.MustCompile(`goacc:"(.*)"`)

func LoadPackage(dirname string) (*packages.Package, error) {

	const mode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
	dir, err := filepath.Abs(dirname)
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

func ParsePackage(pkg *packages.Package) []entity.FileConfig {

	files := []entity.FileConfig{}
	for i := range pkg.GoFiles {
		goFile := pkg.GoFiles[i]
		syntax := pkg.Syntax[i]

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
				namedType, structType, isValid := parseNamedStructType(pkg.TypesInfo.ObjectOf(node.Name))
				if !isValid {
					break
				}
				structConfig := parseStruct(namedType, structType)

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

		files = append(files, *entity.NewFileConfigBuilder(goFile, pkg.Name, imports, structs).Purge())
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

func parseStruct(namedType *types.Named, structType *types.Struct) entity.StructConfig {

	structSupportsBuilder := entity.NewStructSupportsBuilder()

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
		features := parseStructFieldTag([]string{"-"})

		// Parse field tags.
		if parsedTag := regTag.FindStringSubmatch(structType.Tag(i)); len(parsedTag) == 2 {
			features = parseStructFieldTag(strings.Split(parsedTag[1], ","))
		}

		fields = append(fields, *entity.NewFieldConfigBuilder(
			fieldName,
			fieldType.String(),
			&features,
		).Purge())
	}

	return *entity.NewStructConfigBuilder(
		namedType.Obj().Name(),
		*structSupportsBuilder.Purge(),
		"",
		fields,
	).Purge()
}

func parseStructFieldTag(splitedTags []string) entity.FieldConfigFeatures {
	if len(splitedTags) == 1 && splitedTags[0] == "-" {
		return *entity.NewFieldConfigFeaturesBuilder(false, false, false, false, false, false, false).Purge()
	}

	return *entity.NewFieldConfigFeaturesBuilder(
		slices.Contains(splitedTags, "mutex"),
		slices.Contains(splitedTags, "required"),
		slices.Contains(splitedTags, "optional"),
		slices.Contains(splitedTags, "getptr"),
		slices.Contains(splitedTags, "get"),
		slices.Contains(splitedTags, "set"),
		slices.Contains(splitedTags, "override"),
	).Purge()
}

func parseTypeString(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return "*" + parseTypeString(expr.X)
	case *ast.SelectorExpr:
		return parseTypeString(expr.X) + "." + parseTypeString(expr.Sel)
	case *ast.ArrayType:
		len := ""
		if expr.Len != nil {
			len = parseTypeString(expr.Len)
		}
		return "[" + len + "]" + parseTypeString(expr.Elt)
	case *ast.BasicLit:
		return expr.Value

	case *ast.FuncType:
		exprBuilder := strings.Builder{}
		exprBuilder.WriteString("func")

		typeParamExprs := parseTypeStringParam(expr.TypeParams)
		if len(typeParamExprs) > 0 {
			exprBuilder.WriteString("[" + typeParamExprs + "]")
		}
		paramExprs := parseTypeStringParam(expr.Params)
		exprBuilder.WriteString("(" + paramExprs + ")")
		resultExprs := parseTypeStringParam(expr.Results)
		if len(resultExprs) > 0 {
			exprBuilder.WriteString(" (" + resultExprs + ")")
		}
		return exprBuilder.String()

	case *ast.MapType:
		return "map[" + parseTypeString(expr.Key) + "]" + parseTypeString(expr.Value)

	case *ast.InterfaceType:
		return "interface {\n" + parseTypeStringInterfaceParam(expr.Methods) + "\n}"

	case *ast.ChanType:
		switch expr.Dir {
		case ast.RECV & ast.SEND:
			return "chan " + parseTypeString(expr.Value)
		case ast.RECV:
			return "chan <- " + parseTypeString(expr.Value)
		case ast.SEND:
			return "chan " + parseTypeString(expr.Value) + " <-"
		default:
			panic("unknown channel direction")
		}

	case *ast.StructType:
		return "struct {\n" + parseTypeStringStructParam(expr.Fields) + "\n}"

	case *ast.Ellipsis:
		return "..." + parseTypeString(expr.Elt)

	default:
		panic(fmt.Sprintf("unsupported type: %T", expr))
	}
}

func parseTypeStringParam(params *ast.FieldList) string {
	if params == nil {
		return ""
	}

	exprs := make([]string, 0, params.NumFields())
	for _, param := range params.List {
		if len(param.Names) > 0 {
			names := make([]string, 0, len(param.Names))
			for i := range param.Names {
				names = append(names, parseTypeString(param.Names[i]))
			}
			exprs = append(exprs, strings.Join(names, ", ")+" "+parseTypeString(param.Type))
		} else {
			exprs = append(exprs, parseTypeString(param.Type))
		}
	}

	return strings.Join(exprs, ", ")
}

func parseTypeStringInterfaceParam(params *ast.FieldList) string {
	if params == nil {
		return ""
	}

	exprs := make([]string, 0, params.NumFields())
	for _, param := range params.List {
		comment := strings.TrimRight(strings.ReplaceAll(param.Comment.Text(), "\n", "\n// "), "//")
		if !strings.HasSuffix(comment, "\n") {
			comment = "// " + comment + "\n"
		}
		fn := parseTypeString(param.Type)
		if len(param.Names) > 0 && strings.HasPrefix(fn, "func") {
			fn = parseTypeString(param.Names[0]) + strings.TrimPrefix(fn, "func")
		}

		exprs = append(exprs, comment+"\n"+fn)
	}

	return strings.Join(exprs, "\n")
}

func parseTypeStringStructParam(params *ast.FieldList) string {
	if params == nil {
		return ""
	}

	exprs := make([]string, 0, params.NumFields())
	for _, param := range params.List {
		comment := strings.TrimRight(strings.ReplaceAll(param.Comment.Text(), "\n", "\n//"), "//")
		if comment != "" {
			if !strings.HasSuffix(comment, "\n") {
				comment += "\n"
			}
		}

		names := make([]string, 0, len(param.Names))
		for i := range param.Names {
			names = append(names, parseTypeString(param.Names[i]))
		}
		fieldType := parseTypeString(param.Type)
		if param.Tag != nil {
			fieldType += " " + param.Tag.Value
		}

		exprs = append(exprs, comment+strings.Join(names, ", ")+" "+fieldType)
	}

	return strings.Join(exprs, "\n")
}
