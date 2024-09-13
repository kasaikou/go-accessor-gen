package entity

type FieldConfigFeatures struct {
	usesMutex    bool `goacc:"required,get"`
	hasRequired  bool `goacc:"required,get"`
	hasOptional  bool `goacc:"required,get"`
	hasPtrGetter bool `goacc:"required,get"`
	hasGetter    bool `goacc:"required,get"`
	hasSetter    bool `goacc:"required,get"`
}
