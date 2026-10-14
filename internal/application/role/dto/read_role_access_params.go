package dto

type ReadRoleAccessParams struct {
	SetSearch bool
	Search    string
	Order     string
	Offset    int
	Limit     int
	RoleID    int
}
