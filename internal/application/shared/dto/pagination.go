package dto

type Pagination struct {
	Page       int  `json:"page" example:"1"`
	PageSize   int  `json:"page_size" example:"20"`
	TotalItems int  `json:"total_items" example:"100"`
	TotalPages int  `json:"total_pages" example:"5"`
	HasNext    bool `json:"has_next" example:"true"`
	HasPrev    bool `json:"has_prev" example:"false"`
}
