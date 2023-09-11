package domain

type Session struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}
