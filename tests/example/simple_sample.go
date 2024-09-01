package tests

type Sample struct {
	name string `goacc:"required,get"`
	piyo int    `goacc:"optional,get"`
	hoge string `goacc:"optional,get,set"`
}
