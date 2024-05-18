package models

import "github.com/harrydayexe/Omni/internal/snowflake"

// User is a struct that represents a user in the system.
type User struct {
	ID       snowflake.Identifier // ID is a unique identifier for the user.
	Username string               // Username is the user's username.
	Posts    []Post               // Posts is a list of posts the user has made.
}

// NewUser creates a new User with the given ID and username.
func NewUser(id snowflake.Identifier, username string) User {
	return User{
		ID:       id,
		Username: username,
	}
}

// Id returns the ID of the user.
func (u User) Id() uint64 {
	return u.ID.Id()
}

// AddPost adds a post to the user's list of posts.
func (u User) AddPost(post Post) {
	u.Posts = append(u.Posts, post)
}
