package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/infrastructure/validator"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/menu/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Create Menu
// @Description  Adds a new menu item to the system. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu body     presenter.CreateMenuRequest true "Menu data to create"
// @Success      201  {object} response.DataPayload         "Successfully create menu"
// @Failure      400  {object} response.ErrorPayload    "Bad Request (Invalid JSON or Validation Error)"
// @Failure      401  {object} response.ErrorPayload    "Unauthorized"
// @Failure      500  {object} response.ErrorPayload    "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus [post]
func CreateMenu(svc contract.MenuUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.CreateMenuRequest
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c,
				response.ErrorParam{
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

		err = svc.CreateMenuUseCase(c.UserContext(), request, userID)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed create menu",
				})
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create menu",
		})
	}
}


