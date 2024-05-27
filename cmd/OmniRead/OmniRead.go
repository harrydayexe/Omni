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

	ur := NewUserRepo()
	pr := NewPostRepo(*ur)
	cr := NewCommentRepo(*ur, *pr)

	if err := cmd.Run(ctx, api.NewHandler(logger, ur, pr, cr), os.Stdout, os.Args); err != nil {
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

	fmt.Fprintf(os.Stdout, "user id1: %v\n", id1)
	fmt.Fprintf(os.Stdout, "user id2: %v\n", id2)

	return &TestUserRepo{
		users: map[uint64]models.User{
			id1.ToInt(): models.NewUser(id1, "Alice", []snowflake.Snowflake{}),
			id2.ToInt(): models.NewUser(id2, "Bob", []snowflake.Snowflake{}),
		},
	}
}

func (r *TestUserRepo) Read(id snowflake.Snowflake) (*models.User, error) {
	user, ok := r.users[id.ToInt()]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (r *TestUserRepo) Create(entity models.User) error {
	r.users[entity.Id().ToInt()] = entity
	return nil
}

func (r *TestUserRepo) Update(entity models.User) error {
	r.users[entity.Id().ToInt()] = entity
	return nil
}

func (r *TestUserRepo) Delete(entity models.User) error {
	delete(r.users, entity.Id().ToInt())
	return nil
}

type TestPostRepo struct {
	posts map[uint64]models.Post
}

func NewPostRepo(userRepo TestUserRepo) *TestPostRepo {
	snowflakeGenerator := snowflake.NewSnowflakeGenerator(1)
	id1 := snowflakeGenerator.NextID()
	id2 := snowflakeGenerator.NextID()

	fmt.Fprintf(os.Stdout, "post id1: %v\n", id1.ToInt())
	fmt.Fprintf(os.Stdout, "post id2: %v\n", id2.ToInt())

	v := make([]models.User, 0, len(userRepo.users))

	for _, value := range userRepo.users {
		v = append(v, value)
	}

	return &TestPostRepo{
		posts: map[uint64]models.Post{
			id1.ToInt(): models.NewPost(id1, v[0].Id(), "Author 1", time.Now(), "Title 1", "Hello, World!", url.URL{}, 20, make([]snowflake.Snowflake, 0), make([]string, 0)),
			id2.ToInt(): models.NewPost(id2, v[0].Id(), "Author 1", time.Now(), "Title 2", "Hello, World!", url.URL{}, 20, make([]snowflake.Snowflake, 0), make([]string, 0)),
		},
	}
}

func (r *TestPostRepo) Read(id snowflake.Snowflake) (*models.Post, error) {
	fmt.Fprintf(os.Stdout, "post id: %v\n", id)
	post, ok := r.posts[id.ToInt()]
	if !ok {
		return nil, nil
	}
	return &post, nil
}

func (r *TestPostRepo) Create(entity models.Post) error {
	r.posts[entity.Id().ToInt()] = entity
	return nil
}

func (r *TestPostRepo) Update(entity models.Post) error {
	r.posts[entity.Id().ToInt()] = entity
	return nil
}

func (r *TestPostRepo) Delete(entity models.Post) error {
	delete(r.posts, entity.Id().ToInt())
	return nil
}

type TestCommentRepo struct {
	comments map[uint64]models.Comment
}

func NewCommentRepo(userRepo TestUserRepo, postRepo TestPostRepo) TestCommentRepo {
	snowflakeGenerator := snowflake.NewSnowflakeGenerator(1)
	id1 := snowflakeGenerator.NextID()
	id2 := snowflakeGenerator.NextID()

	fmt.Fprintf(os.Stdout, "comment id1: %v\n", id1.ToInt())
	fmt.Fprintf(os.Stdout, "comment id2: %v\n", id2.ToInt())

	v := make([]models.User, 0, len(userRepo.users))
	for _, value := range userRepo.users {
		v = append(v, value)
	}

	return TestCommentRepo{
		comments: map[uint64]models.Comment{
			id1.ToInt(): models.NewComment(id1, v[0].Id(), "Author 1", time.Now(), "Hello, World!"),
			id2.ToInt(): models.NewComment(id2, v[0].Id(), "Author 2", time.Now(), "Hello, World!"),
		},
	}
}

func (r TestCommentRepo) Read(id snowflake.Snowflake) (*models.Comment, error) {
	comment, ok := r.comments[id.ToInt()]
	if !ok {
		return nil, nil
	}
	return &comment, nil
}

func (r TestCommentRepo) Create(entity models.Comment) error {
	r.comments[entity.Id().ToInt()] = entity
	return nil
}

func (r TestCommentRepo) Update(entity models.Comment) error {
	r.comments[entity.Id().ToInt()] = entity
	return nil
}

func (r TestCommentRepo) Delete(entity models.Comment) error {
	delete(r.comments, entity.Id().ToInt())
	return nil
}
