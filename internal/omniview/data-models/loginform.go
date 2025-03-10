package datamodels

type FormPage struct {
	Head   Head
	NavBar NavBar
	Form   Form
}

func NewFormPage(pageType string) FormPage {
	return FormPage{
		Head: Head{Title: "Omni | " + pageType},
		NavBar: NavBar{
			ShouldShowLogin: true,
			IsLoggedIn:      false,
		},
		Form: NewForm(),
	}
}

// Form is a struct that holds the values and errors for forms
type Form struct {
	Values map[string]string
	Errors map[string]string
}

func NewForm() Form {
	return Form{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}
