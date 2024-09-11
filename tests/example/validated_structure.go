package tests

import "errors"

type Validated struct {
	Type string `goacc:"required"`
}

func (v *Validated) goaccPostNewHook() error {
	if v.Type == "" {
		return errors.New("type is empty")
	}
	return nil
}
