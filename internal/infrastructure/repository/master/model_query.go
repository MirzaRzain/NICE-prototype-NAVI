package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) ReadModels(ctx context.Context) ([]model.MasterModel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mm.usecase_id, m.usecase_inference_object_id, mm.created_at, mm.updated_at FROM master_model mm LEFT JOIN usecase_inference_object_model m ON m.model_id=mm.id AND m.deleted_at IS NULL ORDER BY mm.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.MasterModel, 0)
	for rows.Next() {
		var item model.MasterModel
		var mappingID sql.NullInt64
		if err := rows.Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size, &item.UsecaseID, &mappingID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if mappingID.Valid {
			item.UsecaseInferenceObjectID = &mappingID.Int64
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadModelsByUsecaseCode(ctx context.Context, code string) ([]model.MasterModel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mm.usecase_id, m.usecase_inference_object_id, mm.created_at, mm.updated_at FROM master_model mm JOIN master_usecase mu ON mu.id=mm.usecase_id AND mu.deleted_at IS NULL LEFT JOIN usecase_inference_object_model m ON m.model_id=mm.id AND m.deleted_at IS NULL WHERE mu.usecase_code=$1 ORDER BY mm.id DESC`, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.MasterModel, 0)
	for rows.Next() {
		var item model.MasterModel
		var mappingID sql.NullInt64
		if err := rows.Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size, &item.UsecaseID, &mappingID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if mappingID.Valid {
			item.UsecaseInferenceObjectID = &mappingID.Int64
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadModel(ctx context.Context, id int64) (model.MasterModel, error) {
	var item model.MasterModel
	var mappingID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mm.usecase_id, m.usecase_inference_object_id, mm.created_at, mm.updated_at FROM master_model mm LEFT JOIN usecase_inference_object_model m ON m.model_id=mm.id AND m.deleted_at IS NULL WHERE mm.id=$1`, id).Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size, &item.UsecaseID, &mappingID, &item.CreatedAt, &item.UpdatedAt)
	if mappingID.Valid {
		item.UsecaseInferenceObjectID = &mappingID.Int64
	}
	return item, err
}

func (r *repo) ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error) {
	var item model.MasterModel
	var mappingID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mm.usecase_id, m.usecase_inference_object_id, mm.created_at, mm.updated_at FROM master_model mm JOIN master_usecase mu ON mu.id=mm.usecase_id AND mu.deleted_at IS NULL LEFT JOIN usecase_inference_object_model m ON m.model_id=mm.id AND m.deleted_at IS NULL WHERE mu.usecase_code=$1 ORDER BY mm.id DESC LIMIT 1`, code).Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size, &item.UsecaseID, &mappingID, &item.CreatedAt, &item.UpdatedAt)
	if mappingID.Valid {
		item.UsecaseInferenceObjectID = &mappingID.Int64
	}
	return item, err
}

func (r *repo) ReadModelByInferenceObjectID(ctx context.Context, inferenceObjectID int64) (model.MasterModel, error) {
	var item model.MasterModel
	var mappingID sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.version, mm.code, mm.size, mm.usecase_id, m.usecase_inference_object_id, mm.created_at, mm.updated_at FROM master_model mm JOIN usecase_inference_object_model m ON m.model_id=mm.id AND m.deleted_at IS NULL WHERE m.usecase_inference_object_id=$1 ORDER BY mm.id DESC LIMIT 1`, inferenceObjectID).Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Version, &item.Code, &item.Size, &item.UsecaseID, &mappingID, &item.CreatedAt, &item.UpdatedAt)
	if mappingID.Valid {
		item.UsecaseInferenceObjectID = &mappingID.Int64
	}
	return item, err
}

func (r *repo) CreateModel(ctx context.Context, item model.MasterModel) (model.MasterModel, error) {
	if item.UsecaseInferenceObjectID == nil {
		return item, sql.ErrNoRows
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return item, err
	}
	defer tx.Rollback()
	if err = tx.QueryRowContext(ctx, `SELECT usecase_id FROM usecase_inference_object WHERE id=$1 AND deleted_at IS NULL`, *item.UsecaseInferenceObjectID).Scan(&item.UsecaseID); err != nil {
		return item, err
	}
	if err = tx.QueryRowContext(ctx, `INSERT INTO master_model (file_url, name, model_name, version, code, size, usecase_id, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,CURRENT_TIMESTAMP) RETURNING id, created_at`, item.FileURL, item.Name, item.ModelName, item.Version, item.Code, item.Size, item.UsecaseID).Scan(&item.ID, &item.CreatedAt); err != nil {
		return item, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO usecase_inference_object_model (model_id, usecase_inference_object_id, created_at) VALUES ($1,$2,CURRENT_TIMESTAMP)`, item.ID, *item.UsecaseInferenceObjectID); err != nil {
		return item, err
	}
	return item, tx.Commit()
}

func (r *repo) UpdateModel(ctx context.Context, item model.MasterModel) error {
	if item.UsecaseInferenceObjectID == nil {
		return sql.ErrNoRows
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = tx.QueryRowContext(ctx, `SELECT usecase_id FROM usecase_inference_object WHERE id=$1 AND deleted_at IS NULL`, *item.UsecaseInferenceObjectID).Scan(&item.UsecaseID); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `UPDATE master_model SET file_url=$1, name=$2, model_name=$3, version=$4, code=$5, size=$6, usecase_id=$7, updated_at=CURRENT_TIMESTAMP WHERE id=$8`, item.FileURL, item.Name, item.ModelName, item.Version, item.Code, item.Size, item.UsecaseID, item.ID)
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
	mapResult, err := tx.ExecContext(ctx, `UPDATE usecase_inference_object_model SET usecase_inference_object_id=$1, updated_at=CURRENT_TIMESTAMP WHERE model_id=$2 AND deleted_at IS NULL`, *item.UsecaseInferenceObjectID, item.ID)
	if err != nil {
		return err
	}
	mapAffected, err := mapResult.RowsAffected()
	if err != nil {
		return err
	}
	if mapAffected == 0 {
		if _, err = tx.ExecContext(ctx, `INSERT INTO usecase_inference_object_model (model_id, usecase_inference_object_id, created_at) VALUES ($1,$2,CURRENT_TIMESTAMP)`, item.ID, *item.UsecaseInferenceObjectID); err != nil {
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
