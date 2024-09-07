package tests

import "github.com/google/uuid"

type Animal struct {
	id   uuid.UUID
	name string
	kind string
}

func (a *Animal) goaccPreNewHook() {}

func (a *Animal) goaccPostNewHook() {}
