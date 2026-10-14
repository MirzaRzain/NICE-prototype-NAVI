package seeder

import (
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
)

type MasterUsecaseSeed struct {
	ID          int
	Name        string
	Code        string
	Description string
}

var masterUsecases = []MasterUsecaseSeed{
	{ID: 1, Name: "Zone Density", Code: "ZONE_DENSITY", Description: "Calculate people density per zone"},
	{ID: 2, Name: "People Dwell Time", Code: "PEOPLE_DWELL_TIME", Description: "Calculate people dwell time"},
	{ID: 3, Name: "People Footfall", Code: "PEOPLE_FOOTFALL", Description: "Count people footfall and direction"},
	{ID: 4, Name: "Receptionist Analytics", Code: "RECEPTIONIST_ANALYTICS", Description: "Analyze receptionist occupancy and recognition"},
	{ID: 5, Name: "Zone Occupancy", Code: "ZONE_OCCUPANCY", Description: "Calculate zone occupancy"},
	{ID: 6, Name: "Face Recognition", Code: "FACE_RECOGNITION", Description: "Recognize faces"},
}

func masterUsecaseSeeder(s goseeder.Seeder) {
	const statement = `
		INSERT INTO master_usecase
			(id, usecase_name, usecase_code, description, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			usecase_name = EXCLUDED.usecase_name,
			usecase_code = EXCLUDED.usecase_code,
			description = EXCLUDED.description,
			updated_by = EXCLUDED.created_by,
			updated_at = CURRENT_TIMESTAMP;
	`
	for _, item := range masterUsecases {
		if _, err := s.DB.Exec(statement, item.ID, item.Name, item.Code, item.Description, "system"); err != nil {
			log.Fatalf("ERROR execute master usecase seeder: %v", err)
		}
	}

	resetSeq := `SELECT setval(pg_get_serial_sequence('master_usecase', 'id'), COALESCE((SELECT MAX(id) FROM master_usecase), 1))`
	if os.Getenv("DB_DRIVER") == "sqlite" {
		resetSeq = "SELECT 1"
	}
	if _, err := s.DB.Exec(resetSeq); err != nil {
		log.Fatalf("ERROR reset master usecase sequence: %v", err)
	}
}
