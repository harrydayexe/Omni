package models

import (
	"encoding/json"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// User is a struct that represents a user in the system.
type User struct {
	id       snowflake.Snowflake   // Id is a unique identifier for the user.
	Username string                // Username is the user's username.
	Posts    []snowflake.Snowflake // Posts is a list of post id's the user has made.
}

// NewUser creates a new User with the given ID and username.
func NewUser(id snowflake.Snowflake, username string, posts []snowflake.Snowflake) User {
	return User{
		id:       id,
		Username: username,
		Posts:    posts,
	}
}

// Id returns the ID of the user.
func (u User) Id() snowflake.Snowflake {
	return u.id
}

func (u User) MarshalJSON() ([]byte, error) {
	var posts []uint64
	for _, post := range u.Posts {
		posts = append(posts, post.ToInt())
	}

	userAltered := struct {
		Id       uint64   `json:"id"`
		Username string   `json:"username"`
		Posts    []uint64 `json:"posts"`
	}{
		Id:       u.id.ToInt(),
		Username: u.Username,
		Posts:    posts,
	}
	return json.Marshal(userAltered)
}
