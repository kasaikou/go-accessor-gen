package entity

type FieldConfig struct {
	name     string `goacc:"required,get"`
	docText  string `goacc:"get,set"`
	typeName string `goacc:"required,get,set"`

	// If empty, it means no json tag.
	jsonTag  string               `goacc:"required,get"`
	features *FieldConfigFeatures `goacc:"required,get,set"`
}
