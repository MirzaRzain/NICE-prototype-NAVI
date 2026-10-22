package repository

import (
	"context"

	"be-dashboard-nba/internal/application/master/dto"
	"be-dashboard-nba/internal/domain/model"
)

type MasterRepository interface {
	ReadUsecases(ctx context.Context) ([]model.MasterUsecase, error)
	ReadUsecaseList(ctx context.Context) ([]model.UsecaseListItem, error)
	ReadUsecasesWithConfig(ctx context.Context) ([]model.UsecaseWithConfig, error)
	ReadUsecase(ctx context.Context, id int64) (model.MasterUsecase, error)
	CreateUsecase(ctx context.Context, usecase model.MasterUsecase, userID string) (model.MasterUsecase, error)
	UpdateUsecase(ctx context.Context, usecase model.MasterUsecase, userID string) error
	ReadModels(ctx context.Context) ([]model.MasterModel, error)
	ReadModelsByUsecaseCode(ctx context.Context, code string) ([]model.MasterModel, error)
	ReadModel(ctx context.Context, id int64) (model.MasterModel, error)
	ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error)
	ReadModelByInferenceObjectID(ctx context.Context, inferenceObjectID int64) (model.MasterModel, error)
	CreateModel(ctx context.Context, item model.MasterModel) (model.MasterModel, error)
	UpdateModel(ctx context.Context, item model.MasterModel) error
	DeleteModel(ctx context.Context, id int64) error
	SetUsecaseInferenceObjects(ctx context.Context, usecaseID int64, assignments []model.InferenceObjectAssignment, userID string) error
	ReadUsecaseInferenceObjects(ctx context.Context, usecaseID int64) ([]model.InferenceObjectAssignment, error)
	ReadUsecaseInferenceObjectOptions(ctx context.Context, usecaseID int64) ([]model.InferenceObjectOption, error)
	ReadUsecaseInferenceObjectOptionsByCode(ctx context.Context, usecaseCode string) ([]model.InferenceObjectOption, error)
	ReadSources(ctx context.Context, params dto.ReadSourcesParams) ([]model.MasterSource, error)
	CountSources(ctx context.Context, params dto.ReadSourcesParams) (int, error)
	CreateVideo(ctx context.Context, video model.MasterVideo, userID string) (model.MasterVideo, error)
	CreateCamera(ctx context.Context, camera model.MasterCamera, userID string) (model.MasterCamera, error)
	UpdateCamera(ctx context.Context, camera model.MasterCamera, userID string) error
	CreateZone(ctx context.Context, zone model.CameraZone, usecaseIDs []int64, userID string) (model.CameraZone, error)
	UpdateZone(ctx context.Context, zone model.CameraZone, usecaseIDs []int64, userID string) error
	DeleteZone(ctx context.Context, cameraID, zoneID int64) error
	ReadZonesByCamera(ctx context.Context, cameraID int64) ([]model.CameraZone, error)
}
