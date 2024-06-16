package storage

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func createNewCommentRepoForTesting(ctx context.Context, t *testing.T, testDataFile string) (*CommentRepo, func()) {
	t.Parallel()

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

	connURL := mySqlContainer.MustConnectionString(ctx, "parseTime=true")
	db, err := sql.Open("mysql", connURL)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}

	return NewCommentRepo(db, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))), cleanUp
}

func TestReadComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	comment, err := commentRepo.Read(ctx, id)
	if err != nil || comment == nil {
		t.Fatalf("failed to read comment: %s", err)
	}

	expected := models.NewComment(
		id,
		snowflake.ParseId(1796290045997481985),
		snowflake.ParseId(1796290045997481984),
		"johndoe",
		time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
		"Example Comment",
	)
	if *comment != expected {
		t.Fatalf("expected comment to be %v, got %v", expected, comment)
	}
}

func TestCreateComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	idGen := snowflake.NewSnowflakeGenerator(0)
	commentId := idGen.NextID()
	time := time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC)
	newComment := models.NewComment(commentId, snowflake.ParseId(1796290045997481985), snowflake.ParseId(1796290045997481984), "johndoe", time, "Example Comment")

	commentRepo.Create(ctx, newComment)

	readComment, err := commentRepo.Read(ctx, commentId)
	if err != nil || readComment == nil {
		t.Fatalf("failed to read comment. error: %s, comment: %v", err, readComment)
	}

	if *readComment != newComment {
		t.Fatalf("expected comment to be %v, got %v", newComment, readComment)
	}
}

func TestUpdateComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)
	postId := snowflake.ParseId(1796290045997481985)
	authorId := snowflake.ParseId(1796290045997481984)
	time := time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC)
	newComment := models.NewComment(id, postId, authorId, "johndoe", time, "Updated Comment")

	commentRepo.Update(ctx, newComment)

	readComment, err := commentRepo.Read(ctx, newComment.Id())
	if err != nil || readComment == nil {
		t.Fatalf("failed to read comment: %s", err)
	}

	if *readComment != newComment {
		t.Fatalf("expected comment to be %v, got %v", newComment, readComment)
	}
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	commentRepo.Delete(ctx, id)

	readComment, err := commentRepo.Read(ctx, id)
	if err != nil {
		t.Fatalf("an error occurred while reading comment: %s", err)
	}

	if readComment != nil {
		t.Fatalf("expected comment to be nil, got %v", readComment)
	}
}

func TestUnknownCommentTableShouldThrowError(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-comment-table.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	_, err := commentRepo.Read(ctx, id)

	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}
}

func TestUnknownUserTableShouldThrowErrorForComments(t *testing.T) {
	ctx := context.Background()

	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-user-table.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	_, err := commentRepo.Read(ctx, id)

	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}
}
