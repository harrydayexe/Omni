type Post struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	MarkdownUrl string    `json:"markdown_url"`
	Description string    `json:"description"`
}
