package dto

type CreateMenuPermissionParams struct {
	Code       string
	ActionName string
	MenuID     int
	CreatedBy  string
}
