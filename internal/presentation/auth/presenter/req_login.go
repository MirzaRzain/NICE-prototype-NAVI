package presenter

import "be-dashboard-nba/internal/application/utils"

type LoginRequest struct {
	PIN string `json:"pin" validate:"required"`
	// Deprecated compatibility fields; device login only uses PIN.
	Email    string `json:"-"`
	Password string `json:"-"`
}

func (request *LoginRequest) ToSessionPayload(userID string, userAgent, iPAddress string) (
	params SessionPayload,
) {
	params = SessionPayload{
		SessionID: utils.GenerateUUID(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: iPAddress,
	}
	return
}
