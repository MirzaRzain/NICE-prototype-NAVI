package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type Role struct {
	ID          int
	Name        string
	Code        string
	Description string
	CreatedBy   string
}

var roles = []Role{
	{
		ID:          1,
		Name:        "Super Admin",
		Code:        "9012 3456 7890",
		Description: "Role Super Admin, memiliki semua akses",
		CreatedBy:   "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
	{
		ID:          2,
		Name:        "User",
		Code:        "9012 3456 9999",
		Description: "Role User Standar, akses terbatas",
		CreatedBy:   "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
}

func roleSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_role(
		id, name, code, description, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO NOTHING;
	`
	for i := range roles {
		_, err := s.DB.Exec(statement,
			roles[i].ID,
			roles[i].Name,
			roles[i].Code,
			roles[i].Description,
			roles[i].CreatedBy,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute role seeder : %v", err.Error())
		}
	}
}
