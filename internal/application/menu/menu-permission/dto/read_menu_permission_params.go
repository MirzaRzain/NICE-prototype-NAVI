package dto

type ReadMenuPermissionParams struct {
	SetSearch bool
	Search    string
	Order     string
	Limit     int
	Offset    int
	MenuID    int
}
