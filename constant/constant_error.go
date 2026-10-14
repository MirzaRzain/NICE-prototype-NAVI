package constant

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// error.
var (
	ErrFailedParseRequest = fiber.NewError(http.StatusBadRequest, "failed to parse request")

	ErrHeaderTokenNotFound = fiber.NewError(http.StatusUnauthorized, "header authorization not found")
	ErrHeaderTokenInvalid  = fiber.NewError(http.StatusUnauthorized, "invalid header token")
	ErrTokenInvalid        = fiber.NewError(http.StatusUnauthorized, "invalid token")
	ErrTokenMissing        = fiber.NewError(http.StatusUnauthorized, "missing token")
	ErrTokenExpired        = fiber.NewError(http.StatusUnauthorized, "expired token")
	ErrTokenUnauthorized   = fiber.NewError(http.StatusUnauthorized, "unauthorized token")

	// 401
	ErrUserIdNotFound           = fiber.NewError(http.StatusNotFound, "user Id not found")
	ErrDataNotFound             = fiber.NewError(http.StatusNotFound, "data not found")
	ErrMenuIdNotFound           = fiber.NewError(http.StatusNotFound, "menu not found")
	ErrMenuPermissionIdNotFound = fiber.NewError(http.StatusNotFound, "menu permission not found")
	ErrRoleIdNotFound           = fiber.NewError(http.StatusNotFound, "role not found")

	//403
	ErrForbiddenSelfDelete = fiber.NewError(http.StatusForbidden, "anda tidak diperbolehkan menghapus akun Anda sendiri")

	ErrMenuHasChildren = fiber.NewError(http.StatusBadRequest, "tidak dapat menambahkan parent, karena menu masih memilki child")

	ErrUnknownSource = fiber.NewError(http.StatusInternalServerError, "an error occurred, please try again later")
)

// error message.
const (
	ErrMsgValidate      = "There are some errors in your request"
	ErrMsgUnknownSource = "an error occurred, please try again later"
)

// error form field.
var (
	// 400.
	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrAccountNotFound      = errors.New("account not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrIDNull               = errors.New("ID cannot be null")
	ErrWrongEmailOrPassword = errors.New("wrong email or password")

	// 403.
	ErrForbiddenRole       = errors.New("your role is not allowed to access this resource")
	ErrForbiddenPermission = errors.New("your permission is not allowed to access this resource")
)
