package response

import (
	"math"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/internal/presentation/request"
)

type DataPayload struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type paginatePayload struct {
	CurrentPage int     `json:"current_page"`
	PerPage     int     `json:"per_page"`
	TotalPage   float64 `json:"total_page"`
	TotalData   int64   `json:"total_data"`
}

type PaginatePayload struct {
	Code     int             `json:"code"`
	Data     interface{}     `json:"data"`
	Message  string          `json:"message"`
	Paginate paginatePayload `json:"paginate"`
}

func Data(c *fiber.Ctx, res DataPayload) error {
	return c.Status(res.Code).JSON(res)
}

func Paginate(c *fiber.Ctx, paginate request.PaginationPayload, totalData int64, res DataPayload) error {
	return c.Status(res.Code).JSON(PaginatePayload{
		Code:     res.Code,
		Data:     res.Data,
		Message:  res.Message,
		Paginate: toPaginatePayload(paginate.Page, paginate.Limit, totalData),
	})
}

func toPaginatePayload(currentPage, perPage int, totalData int64) paginatePayload {
	totalPage := math.Ceil(float64(totalData) / float64(perPage))

	return paginatePayload{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalPage:   totalPage,
		TotalData:   totalData,
	}
}
