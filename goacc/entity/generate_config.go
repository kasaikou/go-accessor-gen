package entity

type GenerateConfig struct {
	includePattern string `goacc:"required,get,json"`
	defaultLabel   string `goacc:"required,get,json"`
}

func (egc *GenerateConfig) goaccPreNewHook() {
	egc.defaultLabel = "-"
}
