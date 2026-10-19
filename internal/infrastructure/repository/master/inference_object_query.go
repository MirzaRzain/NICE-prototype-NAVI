package repository

import (
	"context"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) ReadUsecaseInferenceObjects(ctx context.Context, usecaseID int64) ([]model.InferenceObjectAssignment, error) {
	const query = `SELECT uio.id, mio.id, mio.inference_id, mio.object_id, mi.inference_name, mi.inference_code, mo.object_name, mo.object_code FROM usecase_inference_object uio JOIN master_inference_object mio ON mio.id=uio.inference_object_id JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL WHERE uio.usecase_id=$1 AND uio.deleted_at IS NULL ORDER BY uio.id`
	rows, err := r.db.QueryContext(ctx, query, usecaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make([]model.InferenceObjectAssignment, 0)
	for rows.Next() {
		var item model.InferenceObjectAssignment
		if err = rows.Scan(&item.ID, &item.InferenceObjectID, &item.InferenceID, &item.ObjectID, &item.InferenceName, &item.InferenceCode, &item.ObjectName, &item.ObjectCode); err != nil {
			return nil, err
		}
		data = append(data, item)
	}
	return data, rows.Err()
}
