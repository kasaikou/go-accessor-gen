package entity

type FileConfig struct {
	// Absolute path of Go file.
	filename string `goacc:"required,get"`
	// Package name.
	packageName string `goacc:"required,get"`
	// Configuration of file importing packages.
	imports []ImportConfig `goacc:"required"`
	// Configuration of file defining structs.
	structs []StructConfig `goacc:"required,get"`
}
