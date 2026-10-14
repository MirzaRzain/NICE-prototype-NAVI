package presenter

import (
	"encoding/json"

	"be-dashboard-nba/internal/application/utils"
	"be-dashboard-nba/internal/domain/model"
)

type ReadSourcesRequest struct{ utils.PaginationPayload }

type UsecaseListResponse struct {
	Code    int                     `json:"code"`
	Data    []model.UsecaseListItem `json:"data"`
	Message string                  `json:"message"`
}

type UsecasesWithConfigResponse struct {
	Code    int                       `json:"code"`
	Data    []model.UsecaseWithConfig `json:"data"`
	Message string                    `json:"message"`
}

type CreateUsecaseRequest struct {
	UsecaseName string  `json:"usecase_name" validate:"required,max=255"`
	UsecaseCode string  `json:"usecase_code" validate:"required,max=255"`
	Description *string `json:"description"`
}

type UpdateUsecaseRequest = CreateUsecaseRequest

type CreateModelRequest struct {
	FileURL   string `json:"file_url"`
	Name      string `json:"name" validate:"required,max=255"`
	ModelName string `json:"model_name" validate:"required,max=255"`
	Code      string `json:"code" validate:"required,max=255"`
	Size      int64  `json:"size" validate:"min=0"`
	UsecaseID int64  `json:"usecase_id" validate:"required,min=1"`
}

type UpdateModelRequest = CreateModelRequest

type InferenceObjectRequest struct {
	InferenceID int64 `json:"inference_id" validate:"required,min=1"`
	ObjectID    int64 `json:"object_id" validate:"required,min=1"`
}

type SetInferenceObjectsRequest struct {
	InferenceObjects []InferenceObjectRequest `json:"inference_objects" validate:"dive"`
}

type CreateVideoRequest struct {
	VideoCode   string  `json:"video_code" validate:"required,max=100"`
	VideoName   string  `json:"video_name" validate:"required,max=255"`
	VideoType   string  `json:"video_type" validate:"required,max=50"`
	SourceVideo string  `json:"source_video" validate:"required"`
	ResultVideo *string `json:"result_video"`
}

type CreateCameraRequest struct {
	AreaLocationID  *int64          `json:"area_location_id"`
	CameraCode      *string         `json:"camera_code" validate:"omitempty,max=50"`
	CameraName      *string         `json:"camera_name" validate:"omitempty,max=255"`
	CameraStatus    *bool           `json:"camera_status"`
	IPAddress       *string         `json:"ip_address" validate:"omitempty,max=45"`
	StreamURL       *string         `json:"stream_url" validate:"omitempty,max=255"`
	AuthType        *string         `json:"auth_type" validate:"omitempty,max=50"`
	Username        *string         `json:"username"`
	Password        *string         `json:"password"`
	InstalledAt     *string         `json:"installed_at"`
	UsecaseConfig   json.RawMessage `json:"usecase_config"`
	Channels        *string         `json:"channels"`
	CameraType      *string         `json:"camera_type"`
	CaptureProtocol *string         `json:"capture_protocol"`
	CaptureMode     *string         `json:"capture_mode" validate:"omitempty,oneof=STREAM SNAPSHOT"`
	RTCURL          *string         `json:"rtc_url"`
	DevicePath      *string         `json:"device_path" validate:"omitempty,max=255"`
}

type UpdateCameraRequest = CreateCameraRequest

type CreateZoneRequest struct {
	ZoneName       string          `json:"zone_name" validate:"required,max=255"`
	ZoneType       string          `json:"zone_type" validate:"required,max=50"`
	ZoneCoordinate json.RawMessage `json:"zone_coordinate" validate:"required"`
	Config         json.RawMessage `json:"config"`
	UsecaseIDs     []int64         `json:"usecase_ids" validate:"dive,min=1"`
}

type UpdateZoneRequest struct {
	ZoneName       string          `json:"zone_name" validate:"required,max=255"`
	ZoneType       string          `json:"zone_type" validate:"required,max=50"`
	ZoneCoordinate json.RawMessage `json:"zone_coordinate" validate:"required"`
	Config         json.RawMessage `json:"config"`
	UsecaseIDs     []int64         `json:"usecase_ids" validate:"dive,min=1"`
}
