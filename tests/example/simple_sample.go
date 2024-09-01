package tests

type Sample struct {
	name string `goacc:"required,get,json"`
	piyo int    `goacc:"optional,get,json(,omitempty)"`
	hoge string `goacc:"optional,get,set"`
}
