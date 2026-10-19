package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/internal/domain/model"
)

const modelSelect = `
	SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size,
	       mio_link.inference_object_id, mio.inference_id, mio.object_id,
	       mi.inference_name, mi.inference_code, mo.object_name, mo.object_code,
	       mm.created_at, mm.updated_at
`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMasterModel(scanner rowScanner) (model.MasterModel, error) {
	var item model.MasterModel
	var inferenceObjectID, inferenceID, objectID sql.NullInt64
	var inferenceName, inferenceCode, objectName, objectCode sql.NullString
	if err := scanner.Scan(
		&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size,
		&inferenceObjectID, &inferenceID, &objectID,
		&inferenceName, &inferenceCode, &objectName, &objectCode,
		&item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return item, err
	}
	if inferenceObjectID.Valid {
		item.InferenceObjectID = &inferenceObjectID.Int64
	}
	if inferenceID.Valid {
		item.InferenceID = &inferenceID.Int64
	}
	if objectID.Valid {
		item.ObjectID = &objectID.Int64
	}
	if inferenceName.Valid {
		item.InferenceName = inferenceName.String
	}
	if inferenceCode.Valid {
		item.InferenceCode = inferenceCode.String
	}
	if objectName.Valid {
		item.ObjectName = objectName.String
	}
	if objectCode.Valid {
		item.ObjectCode = objectCode.String
	}
	return item, nil
}

func (r *repo) ReadModels(ctx context.Context) ([]model.MasterModel, error) {
	rows, err := r.db.QueryContext(ctx, modelSelect+` FROM master_model mm LEFT JOIN model_inference_object mio_link ON mio_link.model_id=mm.id AND mio_link.deleted_at IS NULL LEFT JOIN master_inference_object mio ON mio.id=mio_link.inference_object_id LEFT JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL LEFT JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL ORDER BY mm.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.MasterModel, 0)
	for rows.Next() {
		item, err := scanMasterModel(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadModelsByUsecaseCode(ctx context.Context, code string) ([]model.MasterModel, error) {
	rows, err := r.db.QueryContext(ctx, modelSelect+` FROM master_model mm JOIN model_inference_object mio_link ON mio_link.model_id=mm.id AND mio_link.deleted_at IS NULL JOIN master_inference_object mio ON mio.id=mio_link.inference_object_id JOIN usecase_inference_object uio ON uio.inference_object_id=mio.id AND uio.deleted_at IS NULL JOIN master_usecase mu ON mu.id=uio.usecase_id AND mu.deleted_at IS NULL LEFT JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL LEFT JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL WHERE mu.usecase_code=$1 GROUP BY mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mio_link.inference_object_id, mio.inference_id, mio.object_id, mi.inference_name, mi.inference_code, mo.object_name, mo.object_code, mm.created_at, mm.updated_at ORDER BY mm.id DESC`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.MasterModel, 0)
	for rows.Next() {
		item, err := scanMasterModel(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadModel(ctx context.Context, id int64) (model.MasterModel, error) {
	return scanMasterModel(r.db.QueryRowContext(ctx, modelSelect+` FROM master_model mm LEFT JOIN model_inference_object mio_link ON mio_link.model_id=mm.id AND mio_link.deleted_at IS NULL LEFT JOIN master_inference_object mio ON mio.id=mio_link.inference_object_id LEFT JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL LEFT JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL WHERE mm.id=$1 ORDER BY mio_link.id LIMIT 1`, id))
}

func (r *repo) ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error) {
	return scanMasterModel(r.db.QueryRowContext(ctx, modelSelect+` FROM master_model mm JOIN model_inference_object mio_link ON mio_link.model_id=mm.id AND mio_link.deleted_at IS NULL JOIN master_inference_object mio ON mio.id=mio_link.inference_object_id JOIN usecase_inference_object uio ON uio.inference_object_id=mio.id AND uio.deleted_at IS NULL JOIN master_usecase mu ON mu.id=uio.usecase_id AND mu.deleted_at IS NULL LEFT JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL LEFT JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL WHERE mu.usecase_code=$1 ORDER BY mm.id DESC, mio_link.id LIMIT 1`, code))
}

func (r *repo) ReadModelByInferenceObjectID(ctx context.Context, inferenceObjectID int64) (model.MasterModel, error) {
	return scanMasterModel(r.db.QueryRowContext(ctx, modelSelect+` FROM master_model mm JOIN model_inference_object mio_link ON mio_link.model_id=mm.id AND mio_link.deleted_at IS NULL JOIN master_inference_object mio ON mio.id=mio_link.inference_object_id LEFT JOIN master_inference mi ON mi.id=mio.inference_id AND mi.deleted_at IS NULL LEFT JOIN master_object mo ON mo.id=mio.object_id AND mo.deleted_at IS NULL WHERE mio.id=$1 ORDER BY mm.id DESC LIMIT 1`, inferenceObjectID))
}

func (r *repo) CreateModel(ctx context.Context, item model.MasterModel) (model.MasterModel, error) {
	if len(item.InferenceObjectIDs) == 0 {
		return item, sql.ErrNoRows
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return item, err
	}
	defer tx.Rollback()
	if err = tx.QueryRowContext(ctx, `INSERT INTO master_model (file_url, name, model_name, version, code, size, created_at) VALUES ($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP) RETURNING id, created_at`, item.FileURL, item.Name, item.ModelName, item.Version, item.Code, item.Size).Scan(&item.ID, &item.CreatedAt); err != nil {
		return item, err
	}
	for _, inferenceObjectID := range item.InferenceObjectIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO model_inference_object (model_id, inference_object_id, created_at) VALUES ($1,$2,CURRENT_TIMESTAMP)`, item.ID, inferenceObjectID); err != nil {
			return item, err
		}
	}
	return item, tx.Commit()
}

func (r *repo) UpdateModel(ctx context.Context, item model.MasterModel) error {
	if len(item.InferenceObjectIDs) == 0 {
		return sql.ErrNoRows
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE master_model SET file_url=$1, name=$2, model_name=$3, version=$4, code=$5, size=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`, item.FileURL, item.Name, item.ModelName, item.Version, item.Code, item.Size, item.ID)
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
	if _, err = tx.ExecContext(ctx, `DELETE FROM model_inference_object WHERE model_id=$1`, item.ID); err != nil {
		return err
	}
	for _, inferenceObjectID := range item.InferenceObjectIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO model_inference_object (model_id, inference_object_id, created_at) VALUES ($1,$2,CURRENT_TIMESTAMP)`, item.ID, inferenceObjectID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repo) DeleteModel(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM master_model WHERE id=$1`, id)
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
