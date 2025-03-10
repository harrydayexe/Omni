package datamodels

import "context"

type NavBar struct {
	// ShouldShowLogin is a boolean that determines if the login button should be shown
	ShouldShowLogin bool
	// IsLoggedIn is a boolean that determines if the user is logged in
	IsLoggedIn bool
}

func NewNavBar(ctx context.Context) NavBar {
	return NavBar{
		ShouldShowLogin: true,
		IsLoggedIn:      ctx.Value("is-logged-in").(bool),
	}
}
