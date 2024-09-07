package entity

type GenerateConfig struct {
	includePattern string `goacc:"required,get,json"`
	defaultTag     string `goacc:"required,get,json"`
}

func (egc *GenerateConfig) goaccPreNewHook() {
	egc.defaultTag = "-"
}
