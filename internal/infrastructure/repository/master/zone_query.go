package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) DeleteZone(ctx context.Context, cameraID, zoneID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM master_camera_zone WHERE id=$1 AND camera_id=$2`, zoneID, cameraID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repo) ReadZonesByCamera(ctx context.Context, cameraID int64) ([]model.CameraZone, error) {
	const query = `
		SELECT z.id, z.camera_id, z.zone_name, z.zone_type, z.zone_coordinate, z.config,
		       z.created_at, z.updated_at, uz.usecase_id, mu.usecase_name, mu.usecase_code
		FROM master_camera_zone z
		LEFT JOIN usecase_zone uz ON uz.master_camera_zone_id=z.id AND uz.deleted_at IS NULL
		LEFT JOIN master_usecase mu ON mu.id=uz.usecase_id AND mu.deleted_at IS NULL
		WHERE z.camera_id=$1 AND z.deleted_at IS NULL
		ORDER BY z.id, uz.usecase_id`
	rows, err := r.db.QueryContext(ctx, query, cameraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.CameraZone, 0)
	byID := make(map[int64]int)
	for rows.Next() {
		var zone model.CameraZone
		var coordinate, config []byte
		var updatedAt sql.NullTime
		var usecaseID sql.NullInt64
		var usecaseName, usecaseCode sql.NullString
		if err = rows.Scan(&zone.ID, &zone.CameraID, &zone.ZoneName, &zone.ZoneType, &coordinate, &config, &zone.CreatedAt, &updatedAt, &usecaseID, &usecaseName, &usecaseCode); err != nil {
			return nil, err
		}
		zone.ZoneCoordinate = json.RawMessage(coordinate)
		zone.Config = json.RawMessage(config)
		if updatedAt.Valid {
			zone.UpdatedAt = &updatedAt.Time
		}
		index, exists := byID[zone.ID]
		if !exists {
			zone.Usecases = make([]model.UsecaseSummary, 0)
			result = append(result, zone)
			index = len(result) - 1
			byID[zone.ID] = index
		}
		if usecaseID.Valid {
			result[index].Usecases = append(result[index].Usecases, model.UsecaseSummary{ID: usecaseID.Int64, Name: usecaseName.String, Code: usecaseCode.String})
		}
	}
	return result, rows.Err()
}
