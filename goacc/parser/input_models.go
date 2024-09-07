//go:generate goacc -i input_models.go
package parser

import "golang.org/x/tools/go/packages"

type LoadPackageInput struct {
	dirname string `goacc:"required,json"`
}

type ParsePackageInput struct {
	pkg        *packages.Package `goacc:"required"`
	defaultTag string            `goacc:"required,json"`
}
