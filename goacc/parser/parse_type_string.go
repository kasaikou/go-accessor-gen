package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

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
