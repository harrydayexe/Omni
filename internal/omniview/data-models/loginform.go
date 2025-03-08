package datamodels

// LoginForm is a struct that holds the values and errors for the login form
type LoginForm struct {
	Values map[string]string
	Errors map[string]string
}

func NewLoginForm() LoginForm {
	return LoginForm{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}
