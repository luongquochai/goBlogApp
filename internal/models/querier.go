// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"context"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeletePostByAuthor(ctx context.Context, arg DeletePostByAuthorParams) error
	GetPostByID(ctx context.Context, id int32) (Post, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListPosts(ctx context.Context) ([]ListPostsRow, error)
	UpdatePassword(ctx context.Context, arg UpdatePasswordParams) error
	UpdatePostByAuthor(ctx context.Context, arg UpdatePostByAuthorParams) error
}

var _ Querier = (*Queries)(nil)
