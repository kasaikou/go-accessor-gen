package tests

import "github.com/google/uuid"

type Animal struct {
	id   uuid.UUID `goacc:"required,get"`
	name string    `goacc:"required,optional,get"` // name of animal
	kind string    `goacc:"required,optional,get"` // kind of animal
}

func (a *Animal) goaccPreNewHook() {}

func (a *Animal) goaccPostNewHook() {}
