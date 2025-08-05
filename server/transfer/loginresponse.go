package transfer

// LoginResponse represents the login response payload
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// LoginData contains the login response data
type LoginData struct {
	Auth LoginResponse `json:"auth"`
}
