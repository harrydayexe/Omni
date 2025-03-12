package datamodels

import (
	"github.com/harrydayexe/Omni/internal/storage"
	"github.com/harrydayexe/Omni/internal/utilities"
)

// Post is a struct that represents all the posts in the database
type AllPosts struct {
	Posts      []PostAndUsername `json:"posts"`
	TotalPages int               `json:"total_pages"`
}

type PostAndUsername struct {
	Username string       `json:"username"`
	Post     storage.Post `json:"post"`
}

func NewAllPosts(posts []storage.GetPostsPagedRow) AllPosts {
	var totalPages int = 0
	if len(posts) > 0 {
		totalPages = int(posts[0].TotalPages)
	}
	return AllPosts{
		Posts: utilities.Map(posts, func(row storage.GetPostsPagedRow) PostAndUsername {
			return PostAndUsername{
				Username: row.Username,
				Post:     row.Post,
			}
		}),
		TotalPages: totalPages,
	}
}
