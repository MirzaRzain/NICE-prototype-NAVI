package response

import (
	"database/sql"
	"errors"
	"net/http"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/presentation/validator"

	"be-dashboard-nba/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
)

type ErrorPayload struct {
	Code    int         `json:"code"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message"`
}

type ErrorParam struct {
	Code    int
	Error   error
	Message string
}

func ErrorValidate(c *fiber.Ctx, err error) error {
	logger.WithContext(c.Context()).Error(err, "error validation request")

	return c.Status(http.StatusBadRequest).JSON(ErrorPayload{
		Code:    http.StatusBadRequest,
		Error:   validator.ValidationErrors(err),
		Message: constant.ErrMsgValidate,
	})
}

func Error(c *fiber.Ctx, res ErrorParam) error {
	var errDetail interface{}

	if res.Error != nil && res.Error.Error() != "" {
		errDetail = res.Error.Error()
	}

	return c.Status(res.Code).JSON(ErrorPayload{
		Code:    res.Code,
		Error:   errDetail,
		Message: res.Message,
	})
}

func ErrorField(c *fiber.Ctx, field, message string) error {
	logger.WithContext(c.Context()).Error(nil, "business validation error", "field", field, "message", message)

	return c.Status(http.StatusBadRequest).JSON(ErrorPayload{
		Code: http.StatusBadRequest,
		Error: map[string]string{
			field: message,
		},
		Message: constant.ErrMsgValidate,
	})
}

func OKForErrNoRows(c *fiber.Ctx, res ErrorParam) error {
	if !errors.Is(res.Error, sql.ErrNoRows) {
		return Error(c, res)
	}

	return c.Status(http.StatusOK).JSON(DataPayload{
		Code:    http.StatusOK,
		Message: res.Message,
		Data:    nil,
	})
}

