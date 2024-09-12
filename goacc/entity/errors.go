package entity

import "errors"

var (
	ErrCannotParseGoFile      = errors.New("cannot parse go file")
	ErrInvalidGoaccFormat     = errors.New("invalid goacc format")
	ErrFailedGoFmtCommand     = errors.New("failed go fmt command")
	ErrFailedGoImportsCommand = errors.New("failed go imports command")
)
