package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/infrastructure/validator"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Create Role
// @Description  Creates a new role. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role body      presenter.CreateRoleRequest true "Role data to create"
// @Success      201  {object}  response.DataPayload   "Successfully create role"
// @Failure      400  {object}  response.ErrorPayload  "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object}  response.ErrorPayload "Unauthorized"
// @Failure      500  {object}  response.ErrorPayload "Failed create role"
// @Security     BearerAuth
// @Router       /roles [post]
func CreateRole(svc contract.RoleUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		var request presenter.CreateRoleRequest
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

		err = svc.CreateRoleUseCase(c.UserContext(), request, userID)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed create role",
				})
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create role",
		})
	}
}


