package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) ReadModels(ctx context.Context) ([]model.MasterModel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, file_url, name, model_name, code, size, usecase_id, created_at, updated_at FROM master_model ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.MasterModel, 0)
	for rows.Next() {
		var item model.MasterModel
		if err := rows.Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Code, &item.Size, &item.UsecaseID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadModel(ctx context.Context, id int64) (model.MasterModel, error) {
	var item model.MasterModel
	err := r.db.QueryRowContext(ctx, `SELECT id, file_url, name, model_name, code, size, usecase_id, created_at, updated_at FROM master_model WHERE id=$1`, id).Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Code, &item.Size, &item.UsecaseID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *repo) ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error) {
	var item model.MasterModel
	err := r.db.QueryRowContext(ctx, `SELECT mm.id, mm.file_url, mm.name, mm.model_name, mm.code, mm.size, mm.usecase_id, mm.created_at, mm.updated_at FROM master_model mm JOIN master_usecase mu ON mu.id=mm.usecase_id AND mu.deleted_at IS NULL WHERE mu.usecase_code=$1 ORDER BY mm.id DESC LIMIT 1`, code).Scan(&item.ID, &item.FileURL, &item.Name, &item.ModelName, &item.Code, &item.Size, &item.UsecaseID, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *repo) CreateModel(ctx context.Context, item model.MasterModel) (model.MasterModel, error) {
	err := r.db.QueryRowContext(ctx, `INSERT INTO master_model (file_url, name, model_name, code, size, usecase_id, created_at) VALUES ($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP) RETURNING id, created_at`, item.FileURL, item.Name, item.ModelName, item.Code, item.Size, item.UsecaseID).Scan(&item.ID, &item.CreatedAt)
	return item, err
}

func (r *repo) UpdateModel(ctx context.Context, item model.MasterModel) error {
	result, err := r.db.ExecContext(ctx, `UPDATE master_model SET file_url=$1, name=$2, model_name=$3, code=$4, size=$5, usecase_id=$6, updated_at=CURRENT_TIMESTAMP WHERE id=$7`, item.FileURL, item.Name, item.ModelName, item.Code, item.Size, item.UsecaseID, item.ID)
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
