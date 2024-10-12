package storage

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harrydayexe/Omni/internal/models"
	"github.com/harrydayexe/Omni/internal/snowflake"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const knownUserId uint64 = 1796290045997481984
const knownPostId uint64 = 1796290045997481985
const knownCommentId uint64 = 1796290045997481986

func createNewCommentRepoForTesting(ctx context.Context, t *testing.T, testDataFile string) (*CommentRepo, *sql.DB, func()) {
	t.Parallel()

	mySqlContainer, err := mysql.Run(ctx, "mysql:8.4.0",
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

	return NewCommentRepo(db, slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))), db, cleanUp
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

	var requiredEntityError *RequiredEntityDoesNotExistError
	if !errors.As(err, &requiredEntityError) {
		t.Fatalf("expected RequiredEntityDoesNotExistError to be thrown, got %s", err)
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

	var requiredEntityError *RequiredEntityDoesNotExistError
	if !errors.As(err, &requiredEntityError) {
		t.Fatalf("expected RequiredEntityDoesNotExistError to be thrown, got %s", err)
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

func TestUpdateComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481986)
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

func TestUpdateCommentDoesNotExist(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481700)
	postId := snowflake.ParseId(1796290045997481985)
	authorId := snowflake.ParseId(1796290045997481984)
	time := time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC)
	newComment := models.NewComment(id, postId, authorId, "johndoe", time, "Updated Comment")

	err := commentRepo.Update(ctx, newComment)

	var notFoundError *NotFoundError
	if !errors.As(err, &notFoundError) {
		t.Fatalf("expected NotFoundError to be thrown, got %s", err)
	}
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481986)

	commentRepo.Delete(ctx, id)

	readComment, err := commentRepo.Read(ctx, id)
	if err != nil {
		t.Fatalf("an error occurred while reading comment: %s", err)
	}

	if readComment != nil {
		t.Fatalf("expected comment to be nil, got %v", readComment)
	}
}

func TestDeleteCommentDoesNotExist(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796290045997481700)

	err := commentRepo.Delete(ctx, id)

	var notFoundError *NotFoundError
	if !errors.As(err, &notFoundError) {
		t.Fatalf("expected NotFoundError to be thrown, got %s", err)
	}
}

func TestUnknownCommentTableShouldThrowError(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-comment-table.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	_, err := commentRepo.Read(ctx, id)

	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}
}

func TestUnknownUserTableShouldThrowErrorForComments(t *testing.T) {
	ctx := context.Background()

	commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo-bad-user-table.sql")
	defer cleanUp()

	id := snowflake.ParseId(1796301682498338817)

	_, err := commentRepo.Read(ctx, id)

	if err == nil {
		t.Fatalf("expected error to be thrown, got nil")
	}
}

func TestGetCommentsForPost(t *testing.T) {
	ctx := context.Background()

	id := snowflake.ParseId(knownPostId)
	expectedComments := []models.Comment{
		models.NewComment(
			snowflake.ParseId(knownCommentId),
			id,
			snowflake.ParseId(knownUserId),
			"johndoe",
			time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC),
			"Example Comment",
		),
		models.NewComment(
			snowflake.ParseId(knownCommentId+1),
			id,
			snowflake.ParseId(knownUserId),
			"johndoe",
			time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
			"Example Comment 2",
		),
		models.NewComment(
			snowflake.ParseId(knownCommentId+2),
			id,
			snowflake.ParseId(knownUserId),
			"johndoe",
			time.Date(2024, 4, 5, 20, 0, 0, 0, time.UTC),
			"Example Comment 3",
		),
		models.NewComment(
			snowflake.ParseId(knownCommentId+3),
			id,
			snowflake.ParseId(knownUserId),
			"johndoe",
			time.Date(2024, 4, 6, 0, 0, 0, 0, time.UTC),
			"Example Comment 4",
		),
	}

	tests := []struct {
		name             string
		ts               time.Time
		limit            int
		expectedComments []int
	}{
		{
			name:             "Get all comments",
			ts:               time.UnixMilli(0),
			limit:            10,
			expectedComments: []int{0, 1, 2, 3},
		},
		{
			name:             "Get 1 comment",
			ts:               time.UnixMilli(0),
			limit:            1,
			expectedComments: []int{0},
		},
		{
			name:             "Get 2 comments",
			ts:               time.UnixMilli(0),
			limit:            2,
			expectedComments: []int{0, 1},
		},
		{
			name:             "Get comments after a certain time",
			ts:               time.Date(2024, 4, 5, 10, 0, 0, 0, time.UTC),
			limit:            10,
			expectedComments: []int{2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commentRepo, _, cleanUp := createNewCommentRepoForTesting(ctx, t, "comment-repo.sql")
			defer cleanUp()

			comments, err := commentRepo.GetCommentsForPost(ctx, id, test.ts, test.limit)
			if err != nil {
				t.Fatalf("failed to read comment: %s", err)
			}

			if len(comments) != len(test.expectedComments) {
				t.Fatalf("expected %d comments, got %d", len(test.expectedComments), len(comments))
			}

			for i, expectedIndex := range test.expectedComments {
				if comments[i] != expectedComments[expectedIndex] {
					t.Fatalf("expected comment %d to be %v, got %v", i, expectedComments[expectedIndex], comments[i])
				}

			}
		})
	}
}
