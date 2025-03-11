package datamodels

// NewUserResponse is the response after creating a new user
type NewUserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type NewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
