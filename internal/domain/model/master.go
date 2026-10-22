package model

import (
	"encoding/json"
	"time"

	shareddto "be-dashboard-nba/internal/application/shared/dto"
)

type MasterVideo struct {
	ID          int64      `json:"id"`
	VideoCode   string     `json:"video_code"`
	VideoName   string     `json:"video_name"`
	VideoType   string     `json:"video_type"`
	SourceVideo string     `json:"source_video"`
	ResultVideo *string    `json:"result_video"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type MasterUsecase struct {
	ID          int64      `json:"id"`
	UsecaseName string     `json:"usecase_name"`
	UsecaseCode string     `json:"usecase_code"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type UsecaseListItem struct {
	UsecaseID   int64  `json:"usecase_id"`
	UsecaseCode string `json:"usecase_code"`
	UsecaseName string `json:"usecase_name"`
}

type UsecaseConfigItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UsecaseWithConfig struct {
	UsecaseID     int64               `json:"usecase_id"`
	UsecaseCode   string              `json:"usecase_code"`
	UsecaseName   string              `json:"usecase_name"`
	UsecaseConfig []UsecaseConfigItem `json:"usecase_config"`
}

type MasterModel struct {
	ID                 int64      `json:"id"`
	FileURL            string     `json:"file_url"`
	Name               string     `json:"name"`
	ModelName          string     `json:"model_name"`
	Version            string     `json:"version"`
	Code               string     `json:"code"`
	Size               int64      `json:"size"`
	InferenceObjectID  *int64     `json:"inference_object_id,omitempty"`
	InferenceID        *int64     `json:"inference_id,omitempty"`
	ObjectID           *int64     `json:"object_id,omitempty"`
	InferenceName      string     `json:"inference_name,omitempty"`
	InferenceCode      string     `json:"inference_code,omitempty"`
	ObjectName         string     `json:"object_name,omitempty"`
	ObjectCode         string     `json:"object_code,omitempty"`
	InferenceObjectIDs []int64    `json:"-"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

type InferenceObjectAssignment struct {
	ID                int64  `json:"id,omitempty"`
	InferenceObjectID int64  `json:"inference_object_id"`
	InferenceID       int64  `json:"inference_id"`
	ObjectID          int64  `json:"object_id"`
	InferenceName     string `json:"inference_name,omitempty"`
	InferenceCode     string `json:"inference_code,omitempty"`
	ObjectName        string `json:"object_name,omitempty"`
	ObjectCode        string `json:"object_code,omitempty"`
}

type InferenceObjectOption struct {
	InferenceObjectID   int64  `json:"inference_object_id"`
	InferenceObjectName string `json:"inference_object_name"`
}

type ModelOption struct {
	ModelID   int64  `json:"model_id"`
	Version   string `json:"version"`
	ModelName string `json:"model_name"`
}

type MasterCamera struct {
	ID              int64           `json:"id"`
	AreaLocationID  *int64          `json:"area_location_id,omitempty"`
	CameraCode      *string         `json:"camera_code,omitempty"`
	CameraName      *string         `json:"camera_name,omitempty"`
	CameraStatus    *bool           `json:"camera_status,omitempty"`
	IPAddress       *string         `json:"ip_address,omitempty"`
	StreamURL       *string         `json:"stream_url,omitempty"`
	AuthType        *string         `json:"auth_type,omitempty"`
	Username        *string         `json:"username,omitempty"`
	Password        *string         `json:"-"`
	InstalledAt     *time.Time      `json:"installed_at,omitempty"`
	UsecaseConfig   json.RawMessage `json:"usecase_config"`
	Channels        *string         `json:"channels,omitempty"`
	CameraType      *string         `json:"camera_type,omitempty"`
	CaptureProtocol *string         `json:"capture_protocol,omitempty"`
	CaptureMode     *string         `json:"capture_mode,omitempty"`
	RTCURL          *string         `json:"rtc_url,omitempty"`
	DevicePath      *string         `json:"device_path,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`
}

type CameraZone struct {
	ID             int64            `json:"id"`
	CameraID       int64            `json:"camera_id"`
	ZoneName       string           `json:"zone_name"`
	ZoneType       string           `json:"zone_type"`
	ZoneCoordinate json.RawMessage  `json:"zone_coordinate"`
	Config         json.RawMessage  `json:"config"`
	Usecases       []UsecaseSummary `json:"usecases"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      *time.Time       `json:"updated_at,omitempty"`
}

type UsecaseSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type MasterSource struct {
	ID           int64     `json:"id"`
	SourceType   string    `json:"source_type"`
	Code         *string   `json:"code,omitempty"`
	Name         *string   `json:"name,omitempty"`
	Type         *string   `json:"type,omitempty"`
	SourceVideo  *string   `json:"source_video,omitempty"`
	ResultVideo  *string   `json:"result_video,omitempty"`
	CameraStatus *bool     `json:"camera_status,omitempty"`
	IPAddress    *string   `json:"ip_address,omitempty"`
	StreamURL    *string   `json:"stream_url,omitempty"`
	DevicePath   *string   `json:"device_path,omitempty"`
	CaptureMode  *string   `json:"capture_mode,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type MasterSourcePaginationResponse struct {
	Data       []MasterSource       `json:"data"`
	Pagination shareddto.Pagination `json:"pagination"`
}
