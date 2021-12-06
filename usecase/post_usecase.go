package usecase

import (
	"context"
	"grpcImageboard/domain"
)

type PostUseCase interface {
	GetPosts(ctx context.Context, threadId string, skip int64, limit int64) ([]*domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) (string, error)
	DeletePost(ctx context.Context, id string) (int64, error)
}

type postUseCase struct {
	Repository domain.PostRepository
}

func NewPostUseCase(repo domain.PostRepository) PostUseCase {
	return &postUseCase{repo}
}

func (u *postUseCase) GetPosts(ctx context.Context, threadId string, skip int64, limit int64) ([]*domain.Post, error) {
	return u.Repository.GetPosts(ctx, threadId, skip, limit)
}

func (u *postUseCase) CreatePost(ctx context.Context, post *domain.Post) (string, error) {
	return u.Repository.CreatePost(ctx, post)
}

func (u *postUseCase) DeletePost(ctx context.Context, id string) (int64, error) {
	return u.Repository.DeletePost(ctx, id)
}
