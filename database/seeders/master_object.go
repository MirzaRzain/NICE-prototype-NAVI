package seeder

import (
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
)

type MasterObjectSeed struct {
	ID          int
	Name        string
	Code        string
	Description string
}

var masterObjects = []MasterObjectSeed{
	{ID: 1, Name: "Face", Code: "face", Description: "Face object"},
	{ID: 2, Name: "Person", Code: "person", Description: "Person object"},
}

func masterObjectSeeder(s goseeder.Seeder) {
	const statement = `
		INSERT INTO master_object
			(id, object_name, object_code, description, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			object_name = EXCLUDED.object_name,
			object_code = EXCLUDED.object_code,
			description = EXCLUDED.description,
			updated_by = EXCLUDED.created_by,
			updated_at = CURRENT_TIMESTAMP;
	`
	for _, item := range masterObjects {
		if _, err := s.DB.Exec(statement, item.ID, item.Name, item.Code, item.Description, "system"); err != nil {
			log.Fatalf("ERROR execute master object seeder: %v", err)
		}
	}

	resetSeq := `SELECT setval(pg_get_serial_sequence('master_object', 'id'), COALESCE((SELECT MAX(id) FROM master_object), 1))`
	if os.Getenv("DB_DRIVER") == "sqlite" {
		resetSeq = "SELECT 1"
	}
	if _, err := s.DB.Exec(resetSeq); err != nil {
		log.Fatalf("ERROR reset master object sequence: %v", err)
	}
}
