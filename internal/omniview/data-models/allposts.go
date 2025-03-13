package datamodels

import (
	"github.com/harrydayexe/Omni/internal/omniread/datamodels"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/harrydayexe/Omni/internal/utilities"
)

type PostListItem struct {
	datamodels.PostAndUsername
	IsDeleteable bool
}

// AllPosts contains the data to drive the "posts" template
type AllPosts struct {
	// An error string to display if it exists
	Error string
	// Defines if the page is for all posts or just one user's posts
	IsUserPage bool
	// The collection of posts to render
	Posts []PostListItem
	// Determines whether to show the previous page button
	HasPrevious bool
	// Determines whether to show the next page button
	HasNext bool
	// The page number of the previous page
	PreviousPageNumber int
	// The page number of the next page
	NextPageNumber int
}

// Create an AllPosts struct with the given data.
// Automatically sets the HasNext and HasPrevious vars depending on if a valid prevNum
// and nextNum is set. If positive, the bools are true, otherwise they are false
func NewAllPosts(
	error string,
	userID snowflake.Identifier,
	posts datamodels.AllPosts,
	isUserPage bool,
	prevNum,
	nextNum int,
) AllPosts {
	var hasPrev, hasNext bool = false, false
	var prevNumP, nextNumP int = 0, 0
	if prevNum > 0 {
		hasPrev = true
		prevNumP = prevNum
	}
	if nextNum > 0 && nextNum <= posts.TotalPages {
		hasNext = true
		nextNumP = nextNum
	}
	postsList := utilities.Map(posts.Posts, convertPostAndUsernameToPostListItem(userID))

	return AllPosts{
		Error:              error,
		Posts:              postsList,
		IsUserPage:         isUserPage,
		HasPrevious:        hasPrev,
		HasNext:            hasNext,
		PreviousPageNumber: prevNumP,
		NextPageNumber:     nextNumP,
	}
}

func convertPostAndUsernameToPostListItem(
	userID snowflake.Identifier,
) func(datamodels.PostAndUsername) PostListItem {
	if userID == nil {
		return func(p datamodels.PostAndUsername) PostListItem {
			return PostListItem{
				PostAndUsername: p,
				IsDeleteable:    false,
			}
		}
	}
	return func(p datamodels.PostAndUsername) PostListItem {
		return PostListItem{
			PostAndUsername: p,
			IsDeleteable:    p.Post.UserID == int64(userID.Id().ToInt()),
		}
	}
}
