package entity

type StructSupports struct {
	hasPreNewHook       bool `goacc:"optional,get"`
	hasPostNewHook      bool `goacc:"optional,get"`
	hasPostNewHookError bool `goacc:"optional,get"`
}
