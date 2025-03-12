package datamodels

// User is a data model for the "user" partial template
type User struct {
	// The username of the user
	Username string
	// An AllPosts struct containing information about the user's posts
	AllPosts AllPosts
}
