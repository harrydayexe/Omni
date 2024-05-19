package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/harrydayexe/Omni/internal/cmd"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/omniread/api"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()

	if err := cmd.Run(ctx, api.NewHandler(logger, NewUserRepo(), NewPostRepo()), os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

type TestUserRepo struct {
	users map[uint64]models.User
}

func NewUserRepo() *TestUserRepo {
	snowflakeGenerator := snowflake.NewSnowflakeGenerator(1)
	id1 := snowflakeGenerator.NextID()
	id2 := snowflakeGenerator.NextID()

	fmt.Fprintf(os.Stdout, "user id1: %v\n", id1.Id())
	fmt.Fprintf(os.Stdout, "user id2: %v\n", id2.Id())

	return &TestUserRepo{
		users: map[uint64]models.User{
			id1.Id(): models.NewUser(id1, "Alice"),
			id2.Id(): models.NewUser(id2, "Bob"),
		},
	}
}

func (r *TestUserRepo) Read(id snowflake.Identifier) (models.User, error) {
	user, ok := r.users[id.Id()]
	if !ok {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *TestUserRepo) Create(entity models.User) error {
	r.users[entity.Id()] = entity
	return nil
}

func (r *TestUserRepo) Update(entity models.User) error {
	r.users[entity.Id()] = entity
	return nil
}

func (r *TestUserRepo) Delete(entity models.User) error {
	delete(r.users, entity.Id())
	return nil
}

type TestPostRepo struct {
	posts map[uint64]models.Post
}

func NewPostRepo() *TestPostRepo {
	snowflakeGenerator := snowflake.NewSnowflakeGenerator(1)
	id1 := snowflakeGenerator.NextID()
	id2 := snowflakeGenerator.NextID()

	fmt.Fprintf(os.Stdout, "post id1: %v\n", id1.Id())
	fmt.Fprintf(os.Stdout, "post id2: %v\n", id2.Id())

	userRepo := NewUserRepo()
	v := make([]models.User, 0, len(userRepo.users))

	for _, value := range userRepo.users {
		v = append(v, value)
	}

	return &TestPostRepo{
		posts: map[uint64]models.Post{
			id1.Id(): models.NewPost(id1, v[0], time.Now(), "Title 1", "Hello, World!", url.URL{}, 20, make([]models.Comment, 0), make([]string, 0)),
			id2.Id(): models.NewPost(id2, v[0], time.Now(), "Title 2", "Hello, World!", url.URL{}, 20, make([]models.Comment, 0), make([]string, 0)),
		},
	}
}

func (r *TestPostRepo) Read(id snowflake.Identifier) (models.Post, error) {
	fmt.Fprintf(os.Stdout, "post id: %v\n", id)
	post, ok := r.posts[id.Id()]
	if !ok {
		return models.Post{}, fmt.Errorf("post not found")
	}
	return post, nil
}

func (r *TestPostRepo) Create(entity models.Post) error {
	r.posts[entity.Id()] = entity
	return nil
}

func (r *TestPostRepo) Update(entity models.Post) error {
	r.posts[entity.Id()] = entity
	return nil
}

func (r *TestPostRepo) Delete(entity models.Post) error {
	delete(r.posts, entity.Id())
	return nil
}
