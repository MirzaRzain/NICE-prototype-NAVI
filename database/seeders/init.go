package seeder

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.Register(masterInferenceSeeder)
	goseeder.Register(masterObjectSeeder)
}
