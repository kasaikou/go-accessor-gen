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

type ImportConfig struct {
	name string `goacc:"required"`
	path string `goacc:"required"`
}

type StructConfig struct {
	name              string         `goacc:"required,get"`
	docText           string         `goacc:"optional,get,set"`
	defineFilename    string         `goacc:"optional,get,set"`
	structSupports    StructSupports `goacc:"required,getptr"`
	mutexFieldName    string         `goacc:"required,get"`
	enableMarshalJSON bool           `goacc:"required,get"`
	fields            []FieldConfig  `goacc:"required,get"`
}

type StructSupports struct {
	hasPreNewHook  bool `goacc:"optional,get"`
	hasPostNewHook bool `goacc:"optional,get"`
}

type FieldConfig struct {
	name     string `goacc:"required,get"`
	docText  string `goacc:"get,set"`
	typeName string `goacc:"required,get,set"`

	// If empty, it means no json tag.
	jsonTag  string               `goacc:"required,get"`
	features *FieldConfigFeatures `goacc:"required,get,set"`
}

type FieldConfigFeatures struct {
	usesMutex    bool `goacc:"required,get"`
	hasRequired  bool `goacc:"required,get"`
	hasOptional  bool `goacc:"required,get"`
	hasPtrGetter bool `goacc:"required,get"`
	hasGetter    bool `goacc:"required,get"`
	hasSetter    bool `goacc:"required,get"`
}
