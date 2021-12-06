package domain

import (
	"context"
)

type BoardRepository interface {
	GetBoardsStream(ctx context.Context, changed func([]*Board))
	GetBoards(ctx context.Context, skip int64, limit int64) ([]*Board, error)
	GetBoardById(ctx context.Context, id string) (*Board, error)
	CreateBoard(ctx context.Context, board *Board) (string, error)
	DeleteBoard(ctx context.Context, id string) (int64, error)
}

type ThreadRepository interface {
	GetThreads(ctx context.Context, boardId string, skip int64, limit int64) ([]*Thread, error)
	GetHotThreads(ctx context.Context) ([]*Thread, error)
	GetLatestThreads(ctx context.Context) ([]*Thread, error)
	GetThreadById(ctx context.Context, id string) (*Thread, error)
	CreateThread(ctx context.Context, thread *Thread) (string, error)
	UpdateThread(ctx context.Context, thread *Thread) (*Thread, error)
	DeleteThread(ctx context.Context, id string) (int64, error)
}

type PostRepository interface {
	GetPosts(ctx context.Context, threadId string, skip int64, limit int64) ([]*Post, error)
	CreatePost(ctx context.Context, post *Post) (string, error)
	DeletePost(ctx context.Context, id string) (int64, error)
}
