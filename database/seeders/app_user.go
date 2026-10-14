package seeder

import (
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Name      string
	FullName  string
	Email     string
	Password  string
	Active    bool
	Phone     string
	CreatedBy string
}

var users = []User{
	{
		ID:        "fcdb2142-4731-470e-8a1b-8d0037665fb2",
		Name:      "superadmin",
		FullName:  "Super Admin",
		Email:     "superadmin@gmail.com",
		Password:  "$2a$10$oPZtYt40bTpZxNS9Ayjj1OQizSbZ8f5DcOkjNFBJwkwf3qD8xQLiO",
		Phone:     "123143124",
		Active:    true,
		CreatedBy: "system",
	},
	{
		ID:        "8a92bf0c-9e11-4b3b-9a6b-1b1234567890",
		Name:      "user1",
		FullName:  "User Biasa",
		Email:     "user@demo.com",
		Password:  "$2a$10$oPZtYt40bTpZxNS9Ayjj1OQizSbZ8f5DcOkjNFBJwkwf3qD8xQLiO", // bcrypt: "password"
		Phone:     "08123456789",
		Active:    true,
		CreatedBy: "system",
	},
}

func userSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_user(
		id, name, full_name, email, password, phone, active, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		full_name = EXCLUDED.full_name,
		email = EXCLUDED.email,
		password = EXCLUDED.password,
		phone = EXCLUDED.phone,
		active = EXCLUDED.active,
		updated_by = EXCLUDED.created_by,
		updated_at = CURRENT_TIMESTAMP;
	`

	bootstrapPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if bootstrapPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(bootstrapPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("ERROR hash bootstrap admin password: %v", err)
		}

		users[0].Password = string(hashedPassword)
	}
	if bootstrapEmail := os.Getenv("DEFAULT_ADMIN_EMAIL"); bootstrapEmail != "" {
		users[0].Email = bootstrapEmail
	}

	for i := range users {
		_, err := s.DB.Exec(statement,
			users[i].ID,
			users[i].Name,
			users[i].FullName,
			users[i].Email,
			users[i].Password,
			users[i].Phone,
			users[i].Active,
			users[i].CreatedBy,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute user seeder : %v", err.Error())
		}
	}
}
