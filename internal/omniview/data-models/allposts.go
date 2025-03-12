package datamodels

import (
	"github.com/harrydayexe/Omni/internal/omniread/datamodels"
)

// AllPosts contains the data to drive the "posts" template
type AllPosts struct {
	// An error string to display if it exists
	Error string
	// Defines if the page is for all posts or just one user's posts
	IsUserPage bool
	// The collection of posts to render
	Posts []datamodels.PostAndUsername
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

	return AllPosts{
		Error:              error,
		Posts:              posts.Posts,
		IsUserPage:         isUserPage,
		HasPrevious:        hasPrev,
		HasNext:            hasNext,
		PreviousPageNumber: prevNumP,
		NextPageNumber:     nextNumP,
	}
}
