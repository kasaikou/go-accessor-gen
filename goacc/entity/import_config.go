package entity

type ImportConfig struct {
	name string `goacc:"required"`
	path string `goacc:"required"`
}
