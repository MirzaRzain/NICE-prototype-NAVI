package dto

type ReadListUserParams struct {
	SetSearch bool
	Search    string
	Order     string
	Offset    int
	Limit     int
}
