package repository

import (
	"context"
	"database/sql"

	"be-dashboard-nba/internal/domain/model"
)

func (r *repo) ReadUsecases(ctx context.Context) ([]model.MasterUsecase, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, usecase_name, usecase_code, description, created_at, updated_at FROM master_usecase WHERE deleted_at IS NULL ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.MasterUsecase, 0)
	for rows.Next() {
		var item model.MasterUsecase
		if err := rows.Scan(&item.ID, &item.UsecaseName, &item.UsecaseCode, &item.Description, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repo) ReadUsecase(ctx context.Context, id int64) (model.MasterUsecase, error) {
	var item model.MasterUsecase
	err := r.db.QueryRowContext(ctx, `SELECT id, usecase_name, usecase_code, description, created_at, updated_at FROM master_usecase WHERE id=$1 AND deleted_at IS NULL`, id).Scan(&item.ID, &item.UsecaseName, &item.UsecaseCode, &item.Description, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *repo) CreateUsecase(ctx context.Context, usecase model.MasterUsecase, userID string) (model.MasterUsecase, error) {
	const query = `INSERT INTO master_usecase (usecase_name, usecase_code, description, created_by, created_at) VALUES ($1,$2,$3,$4,CURRENT_TIMESTAMP) RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query, usecase.UsecaseName, usecase.UsecaseCode, usecase.Description, userID).Scan(&usecase.ID, &usecase.CreatedAt)
	return usecase, err
}

func (r *repo) UpdateUsecase(ctx context.Context, usecase model.MasterUsecase, userID string) error {
	const query = `UPDATE master_usecase SET usecase_name=$1, usecase_code=$2, description=$3, updated_by=$4, updated_at=CURRENT_TIMESTAMP WHERE id=$5 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, usecase.UsecaseName, usecase.UsecaseCode, usecase.Description, userID, usecase.ID)
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

func (r *repo) SetUsecaseInferenceObjects(ctx context.Context, usecaseID int64, assignments []model.InferenceObjectAssignment, userID string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	rollback := func(e error) error { _ = tx.Rollback(); return e }
	if _, err = tx.ExecContext(ctx, `DELETE FROM usecase_inference_object WHERE usecase_id=$1`, usecaseID); err != nil {
		return rollback(err)
	}
	const insert = `INSERT INTO usecase_inference_object (usecase_id, inference_id, object_id, created_by, created_at) VALUES ($1,$2,$3,$4,CURRENT_TIMESTAMP)`
	for _, assignment := range assignments {
		if _, err = tx.ExecContext(ctx, insert, usecaseID, assignment.InferenceID, assignment.ObjectID, userID); err != nil {
			return rollback(err)
		}
	}
	return tx.Commit()
}
