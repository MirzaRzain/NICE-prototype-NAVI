package seeder

import (
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
)

type MasterInferenceSeed struct {
	ID          int
	Name        string
	Code        string
	Description string
}

var masterInferences = []MasterInferenceSeed{
	{ID: 1, Name: "Detection", Code: "detection", Description: "Object detection inference"},
	{ID: 2, Name: "Tracking", Code: "tracking", Description: "Object tracking inference"},
	{ID: 3, Name: "Recognition", Code: "recognition", Description: "Object recognition inference"},
}

func masterInferenceSeeder(s goseeder.Seeder) {
	const statement = `
		INSERT INTO master_inference
			(id, inference_name, inference_code, description, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			inference_name = EXCLUDED.inference_name,
			inference_code = EXCLUDED.inference_code,
			description = EXCLUDED.description,
			updated_by = EXCLUDED.created_by,
			updated_at = CURRENT_TIMESTAMP;
	`
	for _, item := range masterInferences {
		if _, err := s.DB.Exec(statement, item.ID, item.Name, item.Code, item.Description, "system"); err != nil {
			log.Fatalf("ERROR execute master inference seeder: %v", err)
		}
	}

	resetSeq := `SELECT setval(pg_get_serial_sequence('master_inference', 'id'), COALESCE((SELECT MAX(id) FROM master_inference), 1))`
	if os.Getenv("DB_DRIVER") == "sqlite" {
		resetSeq = "SELECT 1"
	}
	if _, err := s.DB.Exec(resetSeq); err != nil {
		log.Fatalf("ERROR reset master inference sequence: %v", err)
	}
}
