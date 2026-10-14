package api

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/infrastructure/logger"
)

var errorCodeMap = map[error]int{
	constant.ErrFailedParseRequest:  http.StatusBadRequest,
	constant.ErrHeaderTokenNotFound: http.StatusUnauthorized,
	constant.ErrHeaderTokenInvalid:  http.StatusUnauthorized,
	constant.ErrTokenUnauthorized:   http.StatusUnauthorized,
	constant.ErrTokenInvalid:        http.StatusUnauthorized,
	constant.ErrTokenExpired:        http.StatusUnauthorized,
	constant.ErrForbiddenRole:       http.StatusForbidden,
	constant.ErrForbiddenPermission: http.StatusForbidden,
	constant.ErrDataNotFound:        http.StatusNotFound,

	constant.ErrAccountNotFound:      http.StatusUnauthorized,
	constant.ErrPasswordIncorrect:    http.StatusUnauthorized,
	constant.ErrWrongEmailOrPassword: http.StatusBadRequest,
	constant.ErrUserIdNotFound:       http.StatusNotFound,
}

func errorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		logger.WithContext(c.UserContext()).Error().
			Err(err).
			Str("event", "error_handler").
			Str("method", c.Method()).
			Str("path", c.Path()).
			Msg("request error handled")

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code := fiberErr.Code
			message := strings.ToUpper(fiberErr.Message[:1]) + fiberErr.Message[1:]
			return c.Status(code).JSON(fiber.Map{"code": code, "message": message})
		}

		if errors.Is(err, sql.ErrNoRows) {
			err = constant.ErrDataNotFound
		}

		code := mapErrorCode(err)
		message := mapErrorMessage(err, code)
		if len(message) > 0 {
			message = strings.ToUpper(message[:1]) + message[1:]
		}

		return c.Status(code).JSON(fiber.Map{"code": code, "message": message})
	}
}

func mapErrorCode(err error) int {
	for sentinel, code := range errorCodeMap {
		if errors.Is(err, sentinel) {
			return code
		}
	}
	return http.StatusInternalServerError
}

func mapErrorMessage(err error, code int) string {
	appDebug, _ := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	message := err.Error()

	if !appDebug && code == http.StatusInternalServerError {
		message = constant.ErrUnknownSource.Error()
	}

	return message
}
