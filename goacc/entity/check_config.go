package entity

type CheckConfig struct {
	workingDir     string `goacc:"required,get,json"`
	includePattern string `goacc:"required,get,json"`
}
