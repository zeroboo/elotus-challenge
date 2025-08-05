package services

// Claims represents the JWT claims structure
type Claims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}
