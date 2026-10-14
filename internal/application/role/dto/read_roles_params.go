package dto

type ReadRolesParams struct {
	SetSearch bool
	Search    string
	Order     string
	Offset    int
	Limit     int
}
