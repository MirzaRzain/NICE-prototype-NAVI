package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type UserRole struct {
	UserID string
	RoleID int
}

var userRoles = []UserRole{
	{
		UserID: "fcdb2142-4731-470e-8a1b-8d0037665fb2",
		RoleID: 1, // Super Admin
	},
	{
		UserID: "8a92bf0c-9e11-4b3b-9a6b-1b1234567890",
		RoleID: 2, // User reguler
	},
}

func userRoleSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_user_role(
		user_id, role_id
	)
	VALUES($1, $2)
	ON CONFLICT (user_id, role_id) DO NOTHING;
	`
	for i := range userRoles {
		_, err := s.DB.Exec(statement,
			userRoles[i].UserID,
			userRoles[i].RoleID,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute user_role seeder : %v", err.Error())
		}
	}
}
