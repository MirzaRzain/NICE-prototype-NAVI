package dto

import "time"

type ReadSourcesParams struct {
	SetSearch  bool
	Search     string
	SourceType string
	StartTime  *time.Time
	EndTime    *time.Time
	Offset     int
	Limit      int
}
