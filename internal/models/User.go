package models

import "github.com/harrydayexe/Omni/internal/snowflake"

type User struct {
	ID       snowflake.Identifier // ID is a unique identifier for the user.
	Username string               // Username is the user's username.
	Posts    []Post               // Posts is a list of posts the user has made.
}
