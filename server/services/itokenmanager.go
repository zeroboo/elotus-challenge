package services

// ITokenManager defines the interface for token management services
type ITokenManager interface {
	GenerateToken(userID int, username string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	ExtractTokenFromHeader(authHeader string) string
	HasValidBearerFormat(authHeader string) bool
}
