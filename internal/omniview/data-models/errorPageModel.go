package datamodels

// ErrorPageModel is the data model for the error page template
type ErrorPageModel struct {
	Head  struct{ Title string }
	Error struct {
		Title   string
		Message string
	}
}

// NewErrorPageModel creates a new ErrorPageModel
func NewErrorPageModel(title, message string) ErrorPageModel {
	return ErrorPageModel{
		Head: struct {
			Title string
		}{
			Title: "Omni | " + title,
		},
		Error: struct {
			Title   string
			Message string
		}{
			Title:   title,
			Message: message,
		},
	}
}
