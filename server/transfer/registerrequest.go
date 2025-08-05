package transfer

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
