package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"time"

	"be-dashboard-nba/constant"
	masterDTO "be-dashboard-nba/internal/application/master/dto"
	shareddto "be-dashboard-nba/internal/application/shared/dto"
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/infrastructure/logger"
	"be-dashboard-nba/internal/presentation/master/presenter"
	"github.com/pkg/errors"
)

func (s *useCase) CreateUsecase(ctx context.Context, request presenter.CreateUsecaseRequest, userID string) (model.MasterUsecase, error) {
	data, err := s.newMasterRepo(s.db).CreateUsecase(ctx, model.MasterUsecase{UsecaseName: request.UsecaseName, UsecaseCode: request.UsecaseCode, Description: request.Description}, userID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error creating master usecase")
		return model.MasterUsecase{}, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func (s *useCase) ReadUsecases(ctx context.Context) ([]model.MasterUsecase, error) {
	return s.newMasterRepo(s.db).ReadUsecases(ctx)
}

func (s *useCase) ReadUsecaseList(ctx context.Context) ([]model.UsecaseListItem, error) {
	return s.newMasterRepo(s.db).ReadUsecaseList(ctx)
}

func (s *useCase) ReadUsecasesWithConfig(ctx context.Context) ([]model.UsecaseWithConfig, error) {
	return s.newMasterRepo(s.db).ReadUsecasesWithConfig(ctx)
}

func (s *useCase) ReadUsecase(ctx context.Context, id int64) (model.MasterUsecase, error) {
	return s.newMasterRepo(s.db).ReadUsecase(ctx, id)
}

func (s *useCase) ReadModels(ctx context.Context) ([]model.MasterModel, error) {
	return s.newMasterRepo(s.db).ReadModels(ctx)
}

func (s *useCase) ReadModelsByUsecaseCode(ctx context.Context, code string) ([]model.MasterModel, error) {
	return s.newMasterRepo(s.db).ReadModelsByUsecaseCode(ctx, code)
}

func (s *useCase) ReadModel(ctx context.Context, id int64) (model.MasterModel, error) {
	return s.newMasterRepo(s.db).ReadModel(ctx, id)
}

func (s *useCase) ReadModelByUsecaseCode(ctx context.Context, code string) (model.MasterModel, error) {
	return s.newMasterRepo(s.db).ReadModelByUsecaseCode(ctx, code)
}

func (s *useCase) ReadModelByInferenceObjectID(ctx context.Context, inferenceObjectID int64) (model.MasterModel, error) {
	return s.newMasterRepo(s.db).ReadModelByInferenceObjectID(ctx, inferenceObjectID)
}

func (s *useCase) CreateModel(ctx context.Context, request presenter.CreateModelRequest) (model.MasterModel, error) {
	return s.newMasterRepo(s.db).CreateModel(ctx, model.MasterModel{FileURL: request.FileURL, Name: request.Name, ModelName: request.ModelName, Version: request.Version, Code: request.Code, Size: request.Size, InferenceObjectIDs: request.InferenceObjectIDs})
}

func (s *useCase) UpdateModel(ctx context.Context, id int64, request presenter.UpdateModelRequest) error {
	err := s.newMasterRepo(s.db).UpdateModel(ctx, model.MasterModel{ID: id, FileURL: request.FileURL, Name: request.Name, ModelName: request.ModelName, Version: request.Version, Code: request.Code, Size: request.Size, InferenceObjectIDs: request.InferenceObjectIDs})
	if errors.Is(err, sql.ErrNoRows) {
		return constant.ErrDataNotFound
	}
	return err
}

func (s *useCase) DeleteModel(ctx context.Context, id int64) error {
	err := s.newMasterRepo(s.db).DeleteModel(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return constant.ErrDataNotFound
	}
	return err
}

func (s *useCase) UpdateUsecase(ctx context.Context, usecaseID int64, request presenter.UpdateUsecaseRequest, userID string) error {
	err := s.newMasterRepo(s.db).UpdateUsecase(ctx, model.MasterUsecase{ID: usecaseID, UsecaseName: request.UsecaseName, UsecaseCode: request.UsecaseCode, Description: request.Description}, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ErrDataNotFound
		}
		logger.WithContext(ctx).Error().Err(err).Int64("usecase_id", usecaseID).Msg("error updating master usecase")
		return errors.WithStack(constant.ErrUnknownSource)
	}
	return nil
}

func (s *useCase) SetUsecaseInferenceObjects(ctx context.Context, usecaseID int64, request presenter.SetInferenceObjectsRequest, userID string) error {
	assignments := make([]model.InferenceObjectAssignment, 0, len(request.InferenceObjects))
	for _, item := range request.InferenceObjects {
		assignments = append(assignments, model.InferenceObjectAssignment{InferenceID: item.InferenceID, ObjectID: item.ObjectID})
	}
	if err := s.newMasterRepo(s.db).SetUsecaseInferenceObjects(ctx, usecaseID, assignments, userID); err != nil {
		logger.WithContext(ctx).Error().Err(err).Int64("usecase_id", usecaseID).Msg("error setting usecase inference objects")
		return errors.WithStack(constant.ErrUnknownSource)
	}
	return nil
}

func (s *useCase) ReadUsecaseInferenceObjects(ctx context.Context, usecaseID int64) ([]model.InferenceObjectAssignment, error) {
	data, err := s.newMasterRepo(s.db).ReadUsecaseInferenceObjects(ctx, usecaseID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Int64("usecase_id", usecaseID).Msg("error reading usecase inference objects")
		return nil, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func (s *useCase) ReadUsecaseInferenceObjectOptions(ctx context.Context, usecaseID int64) ([]model.InferenceObjectOption, error) {
	return s.newMasterRepo(s.db).ReadUsecaseInferenceObjectOptions(ctx, usecaseID)
}

func (s *useCase) ReadUsecaseInferenceObjectOptionsByCode(ctx context.Context, usecaseCode string) ([]model.InferenceObjectOption, error) {
	return s.newMasterRepo(s.db).ReadUsecaseInferenceObjectOptionsByCode(ctx, usecaseCode)
}

func (s *useCase) ReadSources(ctx context.Context, params masterDTO.ReadSourcesParams) (model.MasterSourcePaginationResponse, error) {
	repo := s.newMasterRepo(s.db)
	total, err := repo.CountSources(ctx, params)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error counting master sources")
		return model.MasterSourcePaginationResponse{}, errors.WithStack(constant.ErrUnknownSource)
	}
	data, err := repo.ReadSources(ctx, params)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error reading master sources")
		return model.MasterSourcePaginationResponse{}, errors.WithStack(constant.ErrUnknownSource)
	}
	page := 1
	if params.Limit > 0 {
		page = (params.Offset / params.Limit) + 1
	}
	totalPages := 0
	if params.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}
	return model.MasterSourcePaginationResponse{Data: data, Pagination: shareddto.Pagination{Page: page, PageSize: params.Limit, TotalItems: total, TotalPages: totalPages, HasNext: page < totalPages, HasPrev: page > 1}}, nil
}

func (s *useCase) CreateVideo(ctx context.Context, request presenter.CreateVideoRequest, userID string) (model.MasterVideo, error) {
	video := model.MasterVideo{VideoCode: request.VideoCode, VideoName: request.VideoName, VideoType: request.VideoType, SourceVideo: request.SourceVideo, ResultVideo: request.ResultVideo}
	data, err := s.newMasterRepo(s.db).CreateVideo(ctx, video, userID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error creating master video")
		return model.MasterVideo{}, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func (s *useCase) CreateCamera(ctx context.Context, request presenter.CreateCameraRequest, userID string) (model.MasterCamera, error) {
	camera := cameraFromRequest(request)
	if request.InstalledAt != nil && *request.InstalledAt != "" {
		installedAt, parseErr := time.Parse(time.RFC3339, *request.InstalledAt)
		if parseErr != nil {
			return model.MasterCamera{}, constant.ErrFailedParseRequest
		}
		camera.InstalledAt = &installedAt
	}
	data, err := s.newMasterRepo(s.db).CreateCamera(ctx, camera, userID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error creating master camera")
		return model.MasterCamera{}, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func (s *useCase) UpdateCamera(ctx context.Context, cameraID int64, request presenter.UpdateCameraRequest, userID string) error {
	camera := cameraFromRequest(request)
	camera.ID = cameraID
	if request.InstalledAt != nil && *request.InstalledAt != "" {
		installedAt, parseErr := time.Parse(time.RFC3339, *request.InstalledAt)
		if parseErr != nil {
			return constant.ErrFailedParseRequest
		}
		camera.InstalledAt = &installedAt
	}
	if err := s.newMasterRepo(s.db).UpdateCamera(ctx, camera, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ErrDataNotFound
		}
		logger.WithContext(ctx).Error().Err(err).Int64("camera_id", cameraID).Msg("error updating master camera")
		return errors.WithStack(constant.ErrUnknownSource)
	}
	return nil
}

func cameraFromRequest(request presenter.CreateCameraRequest) model.MasterCamera {
	config := request.UsecaseConfig
	if len(config) == 0 {
		config = json.RawMessage(`[]`)
	}
	return model.MasterCamera{AreaLocationID: request.AreaLocationID, CameraCode: request.CameraCode, CameraName: request.CameraName, CameraStatus: request.CameraStatus, IPAddress: request.IPAddress, StreamURL: request.StreamURL, AuthType: request.AuthType, Username: request.Username, Password: request.Password, UsecaseConfig: config, Channels: request.Channels, CameraType: request.CameraType, CaptureProtocol: request.CaptureProtocol, CaptureMode: request.CaptureMode, RTCURL: request.RTCURL, DevicePath: request.DevicePath}
}

func (s *useCase) CreateZone(ctx context.Context, cameraID int64, request presenter.CreateZoneRequest, userID string) (model.CameraZone, error) {
	zone, err := zoneFromRequest(request.ZoneName, request.ZoneType, request.ZoneCoordinate, request.Config)
	if err != nil {
		return model.CameraZone{}, err
	}
	zone.CameraID = cameraID
	data, err := s.newMasterRepo(s.db).CreateZone(ctx, zone, request.UsecaseIDs, userID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Int64("camera_id", cameraID).Msg("error creating camera zone")
		return model.CameraZone{}, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func (s *useCase) UpdateZone(ctx context.Context, cameraID, zoneID int64, request presenter.UpdateZoneRequest, userID string) error {
	zone, err := zoneFromRequest(request.ZoneName, request.ZoneType, request.ZoneCoordinate, request.Config)
	if err != nil {
		return err
	}
	zone.ID, zone.CameraID = zoneID, cameraID
	if err = s.newMasterRepo(s.db).UpdateZone(ctx, zone, request.UsecaseIDs, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ErrDataNotFound
		}
		logger.WithContext(ctx).Error().Err(err).Int64("zone_id", zoneID).Msg("error updating camera zone")
		return errors.WithStack(constant.ErrUnknownSource)
	}
	return nil
}

func (s *useCase) DeleteZone(ctx context.Context, cameraID, zoneID int64) error {
	if err := s.newMasterRepo(s.db).DeleteZone(ctx, cameraID, zoneID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constant.ErrDataNotFound
		}
		logger.WithContext(ctx).Error().Err(err).Int64("zone_id", zoneID).Msg("error deleting camera zone")
		return errors.WithStack(constant.ErrUnknownSource)
	}
	return nil
}

func (s *useCase) ReadZonesByCamera(ctx context.Context, cameraID int64) ([]model.CameraZone, error) {
	data, err := s.newMasterRepo(s.db).ReadZonesByCamera(ctx, cameraID)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Int64("camera_id", cameraID).Msg("error reading camera zones")
		return nil, errors.WithStack(constant.ErrUnknownSource)
	}
	return data, nil
}

func zoneFromRequest(name, zoneType string, coordinate, config json.RawMessage) (model.CameraZone, error) {
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	return model.CameraZone{ZoneName: name, ZoneType: zoneType, ZoneCoordinate: coordinate, Config: config}, nil
}
