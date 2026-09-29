package model

import "time"

type RefreshToken struct {
	ID                  int        `json:"id"`
	UserID              int        `json:"user_id"`
	TokenHash           string     `json:"-"`
	ExpiresAt           time.Time  `json:"expires_at"`
	RevokedAt           *time.Time `json:"revoked_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	ReplacedByTokenHash string     `json:"-"`
}
