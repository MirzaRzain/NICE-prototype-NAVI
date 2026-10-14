package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/infrastructure/validator"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/user/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Summary      Create User
// @Description  Creates a new user with the provided details. Requires Bearer token.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body      presenter.CreateUserRequest true "User data to create"
// @Success      201  {object}  response.DataPayload     "Successfully create user"
// @Failure      400  {object}  response.ErrorPayload "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  response.ErrorPayload "Unauthorized"
// @Failure      404  {object}  response.ErrorPayload "Role ID not found"
// @Failure      500  {object}  response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /users [post]
func CreateUser(svc contract.UserUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.CreateUserRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusBadRequest,
				Message: "Failed parse request",
			})
		}

		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.CreateUserUseCase(c.UserContext(), request, userID)
		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusNotFound,
					Message: constant.ErrRoleIdNotFound.Message,
				})
			}
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed create user",
				})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create user",
		})

	}
}


