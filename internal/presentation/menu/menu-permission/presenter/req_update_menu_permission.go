package presenter

type UpdateMenuPermissionRequest struct {
	Code       string `json:"code" validate:"required,min=1,max=50"`
	ActionName string `json:"action_name" validate:"required,min=1,max=50"`
}
