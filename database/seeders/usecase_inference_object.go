package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type UsecaseInferenceObjectSeed struct {
	UsecaseCode   string
	InferenceCode string
	ObjectCode    string
}

var usecaseInferenceObjects = []UsecaseInferenceObjectSeed{
	{UsecaseCode: "ZONE_DENSITY", InferenceCode: "detection", ObjectCode: "person"},
	{UsecaseCode: "PEOPLE_DWELL_TIME", InferenceCode: "tracking", ObjectCode: "person"},
	{UsecaseCode: "FACE_RECOGNITION", InferenceCode: "detection", ObjectCode: "face"},
}

func usecaseInferenceObjectSeeder(s goseeder.Seeder) {
	const inferenceObjectStatement = `
		INSERT INTO master_inference_object (inference_id, object_id, created_at)
		SELECT i.id, o.id, CURRENT_TIMESTAMP
		FROM master_inference i
		JOIN master_object o ON o.object_code = $2 AND o.deleted_at IS NULL
		WHERE i.inference_code = $1 AND i.deleted_at IS NULL
		ON CONFLICT (inference_id, object_id) DO NOTHING;
	`
	const statement = `
		INSERT INTO usecase_inference_object
			(usecase_id, inference_id, object_id, inference_object_id, created_by, created_at)
		SELECT u.id, i.id, o.id, mio.id, $4, CURRENT_TIMESTAMP
		FROM master_usecase u
		JOIN master_inference i ON i.inference_code = $2 AND i.deleted_at IS NULL
		JOIN master_object o ON o.object_code = $3 AND o.deleted_at IS NULL
		JOIN master_inference_object mio ON mio.inference_id = i.id AND mio.object_id = o.id
		WHERE u.usecase_code = $1 AND u.deleted_at IS NULL
		  AND NOT EXISTS (
			  SELECT 1 FROM usecase_inference_object existing
			  WHERE existing.usecase_id = u.id
			    AND existing.inference_id = i.id
			    AND existing.object_id = o.id
			    AND existing.deleted_at IS NULL
		  );
	`
	for _, item := range usecaseInferenceObjects {
		if _, err := s.DB.Exec(inferenceObjectStatement, item.InferenceCode, item.ObjectCode); err != nil {
			log.Fatalf("ERROR execute master inference object seeder: %v", err)
		}
		if _, err := s.DB.Exec(statement, item.UsecaseCode, item.InferenceCode, item.ObjectCode, "system"); err != nil {
			log.Fatalf("ERROR execute usecase inference object seeder: %v", err)
		}
	}
}
