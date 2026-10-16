package usecase

import (
	"context"

	"be-dashboard-nba/internal/application/master/dto"
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/presentation/master/presenter"
)

type MasterUseCase interface {
	ReadUsecases(ctx context.Context) ([]model.MasterUsecase, error)
	ReadUsecaseList(ctx context.Context) ([]model.UsecaseListItem, error)
	ReadUsecasesWithConfig(ctx context.Context) ([]model.UsecaseWithConfig, error)
	ReadUsecase(ctx context.Context, id int64) (model.MasterUsecase, error)
	CreateUsecase(ctx context.Context, request presenter.CreateUsecaseRequest, userID string) (model.MasterUsecase, error)
	UpdateUsecase(ctx context.Context, usecaseID int64, request presenter.UpdateUsecaseRequest, userID string) error
	ReadModels(ctx context.Context) ([]model.MasterModel, error)
	ReadModelsByUsecaseCode(ctx context.Context, code string) ([]model.MasterModel, error)
	ReadModel(ctx context.Context, id int64) (model.MasterModel, error)
	ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error)
	ReadModelByInferenceObjectID(ctx context.Context, inferenceObjectID int64) (model.MasterModel, error)
	CreateModel(ctx context.Context, request presenter.CreateModelRequest) (model.MasterModel, error)
	UpdateModel(ctx context.Context, id int64, request presenter.UpdateModelRequest) error
	DeleteModel(ctx context.Context, id int64) error
	SetUsecaseInferenceObjects(ctx context.Context, usecaseID int64, request presenter.SetInferenceObjectsRequest, userID string) error
	ReadUsecaseInferenceObjects(ctx context.Context, usecaseID int64) ([]model.InferenceObjectAssignment, error)
	ReadSources(ctx context.Context, params dto.ReadSourcesParams) (model.MasterSourcePaginationResponse, error)
	CreateVideo(ctx context.Context, request presenter.CreateVideoRequest, userID string) (model.MasterVideo, error)
	CreateCamera(ctx context.Context, request presenter.CreateCameraRequest, userID string) (model.MasterCamera, error)
	UpdateCamera(ctx context.Context, cameraID int64, request presenter.UpdateCameraRequest, userID string) error
	CreateZone(ctx context.Context, cameraID int64, request presenter.CreateZoneRequest, userID string) (model.CameraZone, error)
	UpdateZone(ctx context.Context, cameraID, zoneID int64, request presenter.UpdateZoneRequest, userID string) error
	DeleteZone(ctx context.Context, cameraID, zoneID int64) error
	ReadZonesByCamera(ctx context.Context, cameraID int64) ([]model.CameraZone, error)
}
