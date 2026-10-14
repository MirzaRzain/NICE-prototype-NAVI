package seeder

import "github.com/kristijorgji/goseeder"

func init() {
	goseeder.Register(masterUsecaseSeeder)
	goseeder.Register(masterInferenceSeeder)
	goseeder.Register(masterObjectSeeder)
	goseeder.Register(usecaseInferenceObjectSeeder)
	goseeder.Register(usecaseConfigSeeder)
}
