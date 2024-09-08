package entity

type GenerateConfig struct {
	workingDir     string `goacc:"required,get,json"`
	includePattern string `goacc:"required,get,json"`
	defaultTag     string `goacc:"required,get,json"`
}

func (gc *GenerateConfig) goaccPreNewHook() {
	gc.defaultTag = "-"
}
