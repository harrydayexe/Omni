package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const knownUserId uint64 = 1796290045997481984
const knownPostId uint64 = 1796290045997481985
const knownCommentId uint64 = 1796290045997481986

func createNewCommentRepoForTesting(ctx context.Context, t *testing.T, testDataFile string) (*CommentRepo, *sql.DB, func()) {
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

	return NewCommentRepo(db, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))), db, cleanUp
}

func TestReadComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(knownCommentId)

	comment, err := commentRepo.Read(ctx, id)
	if err != nil || comment == nil {
		t.Fatalf("failed to read comment: %s", err)
	}

	expected := models.NewComment(
		id,
		snowflake.ParseId(knownPostId),
		snowflake.ParseId(knownUserId),
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

	commentRepo, db, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	idGen := snowflake.NewSnowflakeGenerator(0)
	commentId := idGen.NextID()
	newTime := time.Date(2024, 4, 4, 11, 4, 3, 0, time.UTC)
	newComment := models.NewComment(commentId, snowflake.ParseId(knownPostId), snowflake.ParseId(knownUserId), "johndoe", newTime, "Example Comment")

	err := commentRepo.Create(ctx, newComment)
	if err != nil {
		t.Fatalf("failed to create comment: %s", err)
	}

	var readId, readPostId, readUserId uint64
	var readContent string
	var readTime time.Time
	if err := db.QueryRow("SELECT id, post_id, user_id, content, created_at FROM Comments WHERE id = ?", commentId.ToInt()).Scan(&readId, &readPostId, &readUserId, &readContent, &readTime); err != nil {
		t.Fatalf("failed to read comment from database: %s", err)
	}

	readComment := models.NewComment(snowflake.ParseId(readId), snowflake.ParseId(readPostId), snowflake.ParseId(readUserId), "johndoe", readTime, readContent)
	if newComment != readComment {
		t.Fatalf("expected comment to be %v, got %v", newComment, readComment)
	}
}

func TestCreateCommentUnknownPost(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	idGen := snowflake.NewSnowflakeGenerator(0)
	commentId := idGen.NextID()
	newTime := time.Date(2024, 4, 4, 11, 4, 3, 0, time.UTC)
	newComment := models.NewComment(commentId, snowflake.ParseId(knownPostId+10), snowflake.ParseId(knownUserId), "johndoe", newTime, "Example Comment")

	err := commentRepo.Create(ctx, newComment)
	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}

	if !strings.HasPrefix(err.Error(), "an unknown database error occurred when creating the comment") {
		t.Fatalf("expected database error to be thrown, got %s", err)
	}
}

func TestCreateCommentUnknownUser(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	idGen := snowflake.NewSnowflakeGenerator(0)
	commentId := idGen.NextID()
	newTime := time.Date(2024, 4, 4, 11, 4, 3, 0, time.UTC)
	newComment := models.NewComment(commentId, snowflake.ParseId(knownPostId), snowflake.ParseId(knownUserId+10), "johndoe", newTime, "Example Comment")

	err := commentRepo.Create(ctx, newComment)
	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}

	if !strings.HasPrefix(err.Error(), "an unknown database error occurred when creating the comment") {
		t.Fatalf("expected database error to be thrown, got %s", err)
	}
}

func TestCreateCommentWithTakenId(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	newTime := time.Date(2024, 4, 4, 11, 4, 3, 0, time.UTC)
	newComment := models.NewComment(snowflake.ParseId(knownCommentId), snowflake.ParseId(knownPostId), snowflake.ParseId(knownUserId), "johndoe", newTime, "Example Comment")

	err := commentRepo.Create(ctx, newComment)
	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}

	var alreadyExistsError *EntityAlreadyExistsError
	if !errors.As(err, &alreadyExistsError) {
		t.Fatalf("expected EntityAlreadyExistsError to be thrown, got %s", err)
	}
}

// func TestUpdateComment(t *testing.T) {
// 	ctx := context.Background()
//
// 	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
// 	defer cleanUp()
//
// 	id := snowflake.ParseId(1796301682498338817)
// 	postId := snowflake.ParseId(1796290045997481985)
// 	authorId := snowflake.ParseId(1796290045997481984)
// 	time := time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC)
// 	newComment := models.NewComment(id, postId, authorId, "johndoe", time, "Updated Comment")
//
// 	commentRepo.Update(ctx, newComment)
//
// 	readComment, err := commentRepo.Read(ctx, newComment.Id())
// 	if err != nil || readComment == nil {
// 		t.Fatalf("failed to read comment: %s", err)
// 	}
//
// 	if *readComment != newComment {
// 		t.Fatalf("expected comment to be %v, got %v", newComment, readComment)
// 	}
// }
//
// func TestDeleteComment(t *testing.T) {
// 	ctx := context.Background()
//
// 	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
// 	defer cleanUp()
//
// 	id := snowflake.ParseId(1796301682498338817)
//
// 	commentRepo.Delete(ctx, id)
//
// 	readComment, err := commentRepo.Read(ctx, id)
// 	if err != nil {
// 		t.Fatalf("an error occurred while reading comment: %s", err)
// 	}
//
// 	if readComment != nil {
// 		t.Fatalf("expected comment to be nil, got %v", readComment)
// 	}
// }
//
// func TestUnknownCommentTableShouldThrowError(t *testing.T) {
// 	ctx := context.Background()
//
// 	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-comment-table.sql")
// 	defer cleanUp()
//
// 	id := snowflake.ParseId(1796301682498338817)
//
// 	_, err := commentRepo.Read(ctx, id)
//
// 	if err == nil {
// 		t.Fatalf("expected error to be thrown, got nil")
// 	}
// }
//
// func TestUnknownUserTableShouldThrowErrorForComments(t *testing.T) {
// 	ctx := context.Background()
//
// 	commentRepo, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-user-table.sql")
// 	defer cleanUp()
//
// 	id := snowflake.ParseId(1796301682498338817)
//
// 	_, err := commentRepo.Read(ctx, id)
//
// 	if err == nil {
// 		t.Fatalf("expected error to be thrown, got nil")
// 	}
// }
