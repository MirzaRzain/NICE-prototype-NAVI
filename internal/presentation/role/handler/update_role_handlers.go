package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/infrastructure/validator"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Update Role
// @Description  Updates an existing role by its ID. Requires Bearer token.
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        role_id  path      int                       true "Role ID"
// @Param        role     body      presenter.UpdateRoleRequest true "Role data to update"
// @Success      200      {object}  response.DataPayload   "Successfully update role"
// @Failure      400      {object}  response.ErrorPayload  "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  response.ErrorPayload "Unauthorized"
// @Failure      404      {object}  response.ErrorPayload "Role ID not found"
// @Failure      500      {object}  response.ErrorPayload "Failed update role"
// @Security     BearerAuth
// @Router       /roles/{role_id} [put]
func UpdateRole(svc contract.RoleUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		roleIDParams := c.Params("role_id")

		roleID, err := strconv.Atoi(roleIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}
		var payload presenter.UpdateRoleRequest
		if err = c.BodyParser(&payload); err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}
		if err := validate.Validate(payload); err != nil {
			return response.ErrorValidate(c, err)
		}

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.UpdateRoleUseCase(c.UserContext(), payload, userID, roleID)
		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrRoleIdNotFound.Message,
					})
			} else {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusInternalServerError,
						Message: "Failed update role",
					})
			}
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update role",
		})
	}
}


