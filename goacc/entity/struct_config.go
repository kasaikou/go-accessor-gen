package entity

type StructConfig struct {
	name              string         `goacc:"required,get"`
	docText           string         `goacc:"optional,get,set"`
	defineFilename    string         `goacc:"optional,get,set"`
	structSupports    StructSupports `goacc:"required,getptr"`
	mutexFieldName    string         `goacc:"required,get"`
	enableMarshalJSON bool           `goacc:"required,get"`
	fields            []FieldConfig  `goacc:"required,get"`
}
