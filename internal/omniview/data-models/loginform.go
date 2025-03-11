package datamodels

import (
	"golang.org/x/net/context"
)

// FormPage is a struct that holds the values for a form page
type FormPage struct {
	Head   Head
	NavBar NavBar
	Form   Form
}

// NewFormPage creates the struct for use on a form page
// This will set the page title to Omni | pageType
func NewFormPage(ctx context.Context, pageType string) FormPage {
	return FormPage{
		Head:   Head{Title: "Omni | " + pageType},
		NavBar: NewNavBar(ctx),
		Form:   NewForm(),
	}
}

// Form is a struct that holds the values and errors for forms
type Form struct {
	Values map[string]string
	Errors map[string]string
}

// NewForm creates a new form struct for use in partial forms
func NewForm() Form {
	return Form{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

// FormSuccess is a struct that holds the values for a successful form submission
type FormSuccess struct {
	// The message to display in the success alert
	Message string
	// The URL to redirect to after the success alert is shown
	RedirectURL string
}

func NewFormSuccess(message, redirectURL string) FormSuccess {
	return FormSuccess{
		Message:     message,
		RedirectURL: redirectURL,
	}
}
