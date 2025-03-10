package datamodels

import "context"

type NavBar struct {
	// ShouldShowLogin is a boolean that determines if the login button should be shown
	ShouldShowLogin bool
	// IsLoggedIn is a boolean that determines if the user is logged in
	IsLoggedIn bool
	// ID is the ID of the user that is logged in
	ID uint64
}

func NewNavBar(ctx context.Context) NavBar {
	var userID uint64
	if ctx.Value("user-id") != nil {
		userID = ctx.Value("user-id").(uint64)
	}
	return NavBar{
		ShouldShowLogin: true,
		IsLoggedIn:      ctx.Value("is-logged-in").(bool),
		ID:              userID,
	}
}
