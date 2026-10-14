package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) CreateVideo(ctx context.Context, video model.MasterVideo, userID string) (model.MasterVideo, error) {
	const query = `INSERT INTO master_video (video_code, video_name, video_type, source_video, result_video, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, video.VideoCode, video.VideoName, video.VideoType, video.SourceVideo, video.ResultVideo, userID).Scan(&video.ID, &video.CreatedAt)
	return video, err
}

func (r *repo) CreateCamera(ctx context.Context, camera model.MasterCamera, userID string) (model.MasterCamera, error) {
	const query = `INSERT INTO master_camera (area_location_id, camera_code, camera_name, camera_status, ip_address, stream_url, auth_type, username, password, installed_at, usecase_config, channels, camera_type, capture_protocol, capture_mode, rtc_url, device_path, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,CURRENT_TIMESTAMP) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, camera.AreaLocationID, camera.CameraCode, camera.CameraName, camera.CameraStatus, camera.IPAddress, camera.StreamURL, camera.AuthType, camera.Username, camera.Password, camera.InstalledAt, camera.UsecaseConfig, camera.Channels, camera.CameraType, camera.CaptureProtocol, camera.CaptureMode, camera.RTCURL, camera.DevicePath, userID).Scan(&camera.ID, &camera.CreatedAt)
	return camera, err
}

func (r *repo) UpdateCamera(ctx context.Context, camera model.MasterCamera, userID string) error {
	const query = `UPDATE master_camera SET area_location_id=$1, camera_code=$2, camera_name=$3, camera_status=$4, ip_address=$5, stream_url=$6, auth_type=$7, username=$8, password=COALESCE($9,password), installed_at=$10, usecase_config=$11, channels=$12, camera_type=$13, capture_protocol=$14, capture_mode=$15, rtc_url=$16, device_path=$17, updated_by=$18, updated_at=CURRENT_TIMESTAMP WHERE id=$19 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, camera.AreaLocationID, camera.CameraCode, camera.CameraName, camera.CameraStatus, camera.IPAddress, camera.StreamURL, camera.AuthType, camera.Username, camera.Password, camera.InstalledAt, camera.UsecaseConfig, camera.Channels, camera.CameraType, camera.CaptureProtocol, camera.CaptureMode, camera.RTCURL, camera.DevicePath, userID, camera.ID)
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

func (r *repo) CreateZone(ctx context.Context, zone model.CameraZone, usecaseIDs []int64, userID string) (model.CameraZone, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return zone, err
	}
	rollback := func(e error) (model.CameraZone, error) { _ = tx.Rollback(); return zone, e }
	const insert = `INSERT INTO master_camera_zone (camera_id, zone_name, zone_type, zone_coordinate, config, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP) RETURNING id, created_at`
	if err = tx.QueryRowContext(ctx, insert, zone.CameraID, zone.ZoneName, zone.ZoneType, zone.ZoneCoordinate, zone.Config, userID).Scan(&zone.ID, &zone.CreatedAt); err != nil {
		return rollback(err)
	}
	if err = assignUsecases(ctx, tx, zone.ID, usecaseIDs, userID); err != nil {
		return rollback(err)
	}
	if err = tx.Commit(); err != nil {
		return zone, err
	}
	return zone, nil
}

func (r *repo) UpdateZone(ctx context.Context, zone model.CameraZone, usecaseIDs []int64, userID string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	rollback := func(e error) error { _ = tx.Rollback(); return e }
	const update = `UPDATE master_camera_zone SET zone_name=$1, zone_type=$2, zone_coordinate=$3, config=$4, updated_by=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$6 AND camera_id=$7 AND deleted_at IS NULL`
	result, err := tx.ExecContext(ctx, update, zone.ZoneName, zone.ZoneType, zone.ZoneCoordinate, zone.Config, userID, zone.ID, zone.CameraID)
	if err != nil {
		return rollback(err)
	}
	if affected, e := result.RowsAffected(); e != nil {
		return rollback(e)
	} else if affected == 0 {
		return rollback(sql.ErrNoRows)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM usecase_zone WHERE master_camera_zone_id=$1`, zone.ID); err != nil {
		return rollback(err)
	}
	if err = assignUsecases(ctx, tx, zone.ID, usecaseIDs, userID); err != nil {
		return rollback(err)
	}
	return tx.Commit()
}

func assignUsecases(ctx context.Context, tx *sql.Tx, zoneID int64, usecaseIDs []int64, userID string) error {
	const insert = `INSERT INTO usecase_zone (usecase_id, master_camera_zone_id, created_by, created_at) VALUES ($1,$2,$3,CURRENT_TIMESTAMP)`
	for _, usecaseID := range usecaseIDs {
		if _, err := tx.ExecContext(ctx, insert, usecaseID, zoneID, userID); err != nil {
			return err
		}
	}
	return nil
}
