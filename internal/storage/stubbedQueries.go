package storage

import "context"

type StubbedQueries struct {
	CreateCommentFn                    func(ctx context.Context, arg CreateCommentParams) error
	CreatePostFn                       func(ctx context.Context, arg CreatePostParams) error
	CreateUserFn                       func(ctx context.Context, arg CreateUserParams) error
	DeleteCommentFn                    func(ctx context.Context, id int64) error
	DeletePostFn                       func(ctx context.Context, id int64) error
	DeleteUserFn                       func(ctx context.Context, id int64) error
	FindCommentAndUserByIDFn           func(ctx context.Context, id int64) (FindCommentAndUserByIDRow, error)
	FindCommentsAndUserByPostIDPagedFn func(ctx context.Context, arg FindCommentsAndUserByPostIDPagedParams) ([]FindCommentsAndUserByPostIDPagedRow, error)
	FindPostByIDFn                     func(ctx context.Context, id int64) (Post, error)
	GetPasswordByIDFn                  func(ctx context.Context, id int64) (string, error)
	GetPostsPagedFn                    func(ctx context.Context, offset int32) ([]GetPostsPagedRow, error)
	GetUserAndPostsByIDPagedFn         func(ctx context.Context, arg GetUserAndPostsByIDPagedParams) ([]GetUserAndPostsByIDPagedRow, error)
	GetUserByIDFn                      func(ctx context.Context, id int64) (GetUserByIDRow, error)
	GetUserByUsernameFn                func(ctx context.Context, username string) (int64, error)
	UpdateCommentFn                    func(ctx context.Context, arg UpdateCommentParams) error
	UpdatePostFn                       func(ctx context.Context, arg UpdatePostParams) error
	UpdateUserFn                       func(ctx context.Context, arg UpdateUserParams) error
}

func (q *StubbedQueries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	return q.CreateCommentFn(ctx, arg)
}

func (q *StubbedQueries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	return q.CreatePostFn(ctx, arg)
}

func (q *StubbedQueries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	return q.CreateUserFn(ctx, arg)
}

func (q *StubbedQueries) DeleteComment(ctx context.Context, id int64) error {
	return q.DeleteCommentFn(ctx, id)
}

func (q *StubbedQueries) DeletePost(ctx context.Context, id int64) error {
	return q.DeletePostFn(ctx, id)
}

func (q *StubbedQueries) DeleteUser(ctx context.Context, id int64) error {
	return q.DeleteUserFn(ctx, id)
}

func (q *StubbedQueries) FindCommentAndUserByID(ctx context.Context, id int64) (FindCommentAndUserByIDRow, error) {
	return q.FindCommentAndUserByIDFn(ctx, id)
}

func (q *StubbedQueries) FindCommentsAndUserByPostIDPaged(ctx context.Context, arg FindCommentsAndUserByPostIDPagedParams) ([]FindCommentsAndUserByPostIDPagedRow, error) {
	return q.FindCommentsAndUserByPostIDPagedFn(ctx, arg)
}

func (q *StubbedQueries) FindPostByID(ctx context.Context, id int64) (Post, error) {
	return q.FindPostByIDFn(ctx, id)
}

func (q *StubbedQueries) GetPasswordByID(ctx context.Context, id int64) (string, error) {
	return q.GetPasswordByIDFn(ctx, id)
}

func (q *StubbedQueries) GetPostsPaged(ctx context.Context, offset int32) ([]GetPostsPagedRow, error) {
	return q.GetPostsPagedFn(ctx, offset)
}

func (q *StubbedQueries) GetUserAndPostsByIDPaged(ctx context.Context, arg GetUserAndPostsByIDPagedParams) ([]GetUserAndPostsByIDPagedRow, error) {
	return q.GetUserAndPostsByIDPagedFn(ctx, arg)
}

func (q *StubbedQueries) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	return q.GetUserByIDFn(ctx, id)
}

func (q *StubbedQueries) GetUserByUsername(ctx context.Context, username string) (int64, error) {
	return q.GetUserByUsernameFn(ctx, username)
}

func (q *StubbedQueries) UpdateComment(ctx context.Context, arg UpdateCommentParams) error {
	return q.UpdateCommentFn(ctx, arg)
}

func (q *StubbedQueries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	return q.UpdatePostFn(ctx, arg)
}

func (q *StubbedQueries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	return q.UpdateUserFn(ctx, arg)
}
