package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type UsecaseConfigSeed struct {
	UsecaseCode string
	Key         string
	Value       string
}

var usecaseConfigs = []UsecaseConfigSeed{
	{UsecaseCode: "ZONE_DENSITY", Key: "max_person_count", Value: "6"},
	{UsecaseCode: "ZONE_DENSITY", Key: "min_count_change_time", Value: "5"},
	{UsecaseCode: "ZONE_DENSITY", Key: "min_status_hold_time", Value: "15"},
	{UsecaseCode: "PEOPLE_DWELL_TIME", Key: "min_dwell_time", Value: "15"},
	{UsecaseCode: "PEOPLE_DWELL_TIME", Key: "person_time_limit", Value: "360"},
	{UsecaseCode: "PEOPLE_FOOTFALL", Key: "direction_type", Value: ""},
	{UsecaseCode: "RECEPTIONIST_ANALYTICS", Key: "max_retry_recog", Value: "3"},
	{UsecaseCode: "RECEPTIONIST_ANALYTICS", Key: "min_count_change", Value: "0"},
	{UsecaseCode: "RECEPTIONIST_ANALYTICS", Key: "min_count_change_time", Value: "0"},
	{UsecaseCode: "ZONE_OCCUPANCY", Key: "max_retry_recog", Value: "3"},
	{UsecaseCode: "ZONE_OCCUPANCY", Key: "min_count_change", Value: "0"},
	{UsecaseCode: "ZONE_OCCUPANCY", Key: "min_count_change_time", Value: "0"},
	{UsecaseCode: "FACE_RECOGNITION", Key: "face_recognition_threshold", Value: "0"},
	{UsecaseCode: "FACE_RECOGNITION", Key: "max_retry_recog", Value: "3"},
}

func usecaseConfigSeeder(s goseeder.Seeder) {
	const updateStatement = `
		UPDATE usecase_config
		SET value = $1, updated_by = $2, updated_at = CURRENT_TIMESTAMP
		WHERE usecase_id = (SELECT id FROM master_usecase WHERE usecase_code = $3 AND deleted_at IS NULL)
		  AND key = $4;
	`
	const insertStatement = `
		INSERT INTO usecase_config (usecase_id, key, value, created_by, created_at)
		SELECT id, $1, $2, $3, CURRENT_TIMESTAMP
		FROM master_usecase
		WHERE usecase_code = $4 AND deleted_at IS NULL
		  AND NOT EXISTS (
			  SELECT 1 FROM usecase_config
			  WHERE usecase_id = master_usecase.id AND key = $1
		  );
	`

	for _, item := range usecaseConfigs {
		if _, err := s.DB.Exec(updateStatement, item.Value, "system", item.UsecaseCode, item.Key); err != nil {
			log.Fatalf("ERROR update usecase config seeder: %v", err)
		}
		if _, err := s.DB.Exec(insertStatement, item.Key, item.Value, "system", item.UsecaseCode); err != nil {
			log.Fatalf("ERROR insert usecase config seeder: %v", err)
		}
	}
}
