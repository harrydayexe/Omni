package storage

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func CreateNewUserRepoForTesting(ctx context.Context, t *testing.T, testDataFile string) (*UserRepo, func()) {
	mySqlContainer, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.4.0"),
		mysql.WithDatabase("omni"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("..", "..", "testdata", testDataFile)),
	)

	if err != nil {
		t.Fatalf("Could not start mysql container: %s", err)
	}
	// Clean up container
	var cleanUp = func() {
		if err := mySqlContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to stop mysql container: %s", err)
		}
	}

	connURL := mySqlContainer.MustConnectionString(ctx)
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}

	return NewUserRepo(db, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))), cleanUp
}

func TestReadUser_NoPosts(t *testing.T) {
	ctx := context.Background()

	userRepo, cleanUp := CreateNewUserRepoForTesting(ctx, t, "user-repo-no-posts.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481984)

	user, err := userRepo.Read(ctx, id)
	if err != nil {
		t.Fatalf("failed to read user: %s", err)
	}

	if user.Id() != id {
		t.Fatalf("expected user id to be %v, got %v", id, user.Id())
	}

	if user.Username != "johndoe" {
		t.Fatalf("expected username to be 'johndoe', got %s", user.Username)
	}
}

func TestReadUser_WithPosts(t *testing.T) {
	ctx := context.Background()

	userRepo, cleanUp := CreateNewUserRepoForTesting(ctx, t, "user-repo-with-posts.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481984)

	user, err := userRepo.Read(ctx, id)
	if err != nil {
		t.Fatalf("failed to read user: %s", err)
	}

	if user.Id() != id {
		t.Fatalf("expected user id to be %v, got %v", id, user.Id())
	}

	if user.Username != "johndoe" {
		t.Fatalf("expected username to be 'johndoe', got %s", user.Username)
	}

	if len(user.Posts) != 2 {
		t.Fatalf("expected user to have 2 posts, got %d", len(user.Posts))
	}

	if user.Posts[0].ToInt() != 1796301682498338816 {
		t.Fatalf("expected first post id to be 1796301682498338816, got %d", user.Posts[0].ToInt())
	}
	if user.Posts[1].ToInt() != 1796301682498338817 {
		t.Fatalf("expected first post id to be 1796301682498338817, got %d", user.Posts[1].ToInt())
	}
}
