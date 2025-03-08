package auth

// LoginResponse is the response from the /api/login endpoint
type LoginResponse struct {
	Token   string `json:"access_token"`
	Type    string `json:"token_type"`
	Expires int    `json:"expires_in"`
}

// LoginRequest is the body data for a request to /api/login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
