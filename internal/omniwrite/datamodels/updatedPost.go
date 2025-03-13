package datamodels

type UpdatedPost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	MarkdownUrl string `json:"markdown_url"`
}
