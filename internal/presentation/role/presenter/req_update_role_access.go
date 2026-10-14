package presenter

type UpdateRoleAccessItem struct {
	AccessID  int   `json:"access_id" validate:"required"`
	HasAccess *bool `json:"has_access" validate:"required"`
}

type UpdateRoleAccessRequest struct {
	AccessItem []UpdateRoleAccessItem `json:"access_item" validate:"required,min=1,dive"`
}
