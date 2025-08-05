package transfer

// RegisterData contains the user data for registration response
type RegisterData struct {
	User UserInfo `json:"user"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
