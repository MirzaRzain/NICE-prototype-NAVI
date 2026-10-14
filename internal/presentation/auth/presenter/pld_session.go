package presenter

import "be-dashboard-nba/internal/application/jwt"

type SessionPayload struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}

func (request *SessionPayload) ToAccessTokenRequest() (
	params jwt.AccessTokenPayload,
) {
	params = jwt.AccessTokenPayload{
		SessionID: request.SessionID,
		UserID:    request.UserID,
	}
	return
}

func (request *SessionPayload) ToRefreshTokenRequest() (
	params jwt.RefreshTokenPayload,
) {
	params = jwt.RefreshTokenPayload{
		SessionID: request.SessionID,
	}
	return
}
