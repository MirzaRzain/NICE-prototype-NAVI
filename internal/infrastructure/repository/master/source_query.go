package repository

import (
	"be-dashboard-nba/internal/application/master/dto"
	"be-dashboard-nba/internal/domain/model"
	"context"
)

const sourceUnion = `
	SELECT id, 'camera' AS source_type, camera_code AS code, camera_name AS name,
	       'camera' AS type, NULL AS source_video, NULL AS result_video,
	       camera_status, ip_address, stream_url, device_path, capture_mode, created_at
	FROM master_camera WHERE deleted_at IS NULL
	UNION ALL
	SELECT id, 'video' AS source_type, video_code AS code, video_name AS name,
	       video_type AS type, source_video, result_video,
	       NULL AS camera_status, NULL AS ip_address, NULL AS stream_url, NULL AS device_path, NULL AS capture_mode, created_at
	FROM master_video WHERE deleted_at IS NULL`

func (r *repo) ReadSources(ctx context.Context, params dto.ReadSourcesParams) ([]model.MasterSource, error) {
	query := `SELECT id, source_type, code, name, type, source_video, result_video, camera_status, ip_address, stream_url, device_path, capture_mode, created_at FROM (` + sourceUnion + `) sources WHERE ($1 = false OR LOWER(code) LIKE LOWER($2) OR LOWER(name) LIKE LOWER($2)) AND ($3 = '' OR source_type = $3) AND ($4 IS NULL OR created_at >= $4) AND ($5 IS NULL OR created_at < $5) ORDER BY created_at DESC, id DESC LIMIT $6 OFFSET $7`
	rows, err := r.db.QueryContext(ctx, query, params.SetSearch, params.Search, params.SourceType, params.StartTime, params.EndTime, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]model.MasterSource, 0)
	for rows.Next() {
		var item model.MasterSource
		if err = rows.Scan(&item.ID, &item.SourceType, &item.Code, &item.Name, &item.Type, &item.SourceVideo, &item.ResultVideo, &item.CameraStatus, &item.IPAddress, &item.StreamURL, &item.DevicePath, &item.CaptureMode, &item.CreatedAt); err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, rows.Err()
}

func (r *repo) CountSources(ctx context.Context, params dto.ReadSourcesParams) (int, error) {
	query := `SELECT COUNT(*) FROM (` + sourceUnion + `) sources WHERE ($1 = false OR LOWER(code) LIKE LOWER($2) OR LOWER(name) LIKE LOWER($2)) AND ($3 = '' OR source_type = $3) AND ($4 IS NULL OR created_at >= $4) AND ($5 IS NULL OR created_at < $5)`
	var count int
	err := r.db.QueryRowContext(ctx, query, params.SetSearch, params.Search, params.SourceType, params.StartTime, params.EndTime).Scan(&count)
	return count, err
}
