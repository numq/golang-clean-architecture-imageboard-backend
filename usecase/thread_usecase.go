package usecase

import (
	"context"
	"grpcImageboard/domain"
)

type ThreadUseCase interface {
	GetThreads(ctx context.Context, boardId string, skip int64, limit int64) ([]*domain.Thread, error)
	GetHotThreads(ctx context.Context) ([]*domain.Thread, error)
	GetLatestThreads(ctx context.Context) ([]*domain.Thread, error)
	GetThreadById(ctx context.Context, id string) (*domain.Thread, error)
	CreateThread(ctx context.Context, thread *domain.Thread) (string, error)
	UpdateThread(ctx context.Context, thread *domain.Thread) (*domain.Thread, error)
	DeleteThread(ctx context.Context, id string) (int64, error)
}

type threadUseCase struct {
	Repository domain.ThreadRepository
}

func NewThreadUseCase(repo domain.ThreadRepository) ThreadUseCase {
	return &threadUseCase{repo}
}

func (u *threadUseCase) GetThreads(ctx context.Context, boardId string, skip int64, limit int64) ([]*domain.Thread, error) {
	return u.Repository.GetThreads(ctx, boardId, skip, limit)
}

func (u *threadUseCase) GetHotThreads(ctx context.Context) ([]*domain.Thread, error) {
	return u.Repository.GetHotThreads(ctx)
}

func (u *threadUseCase) GetLatestThreads(ctx context.Context) ([]*domain.Thread, error) {
	return u.Repository.GetLatestThreads(ctx)
}

func (u *threadUseCase) GetThreadById(ctx context.Context, id string) (*domain.Thread, error) {
	return u.Repository.GetThreadById(ctx, id)
}

func (u *threadUseCase) CreateThread(ctx context.Context, thread *domain.Thread) (string, error) {
	return u.Repository.CreateThread(ctx, thread)
}

func (u *threadUseCase) UpdateThread(ctx context.Context, thread *domain.Thread) (*domain.Thread, error) {
	return u.Repository.UpdateThread(ctx, thread)
}

func (u *threadUseCase) DeleteThread(ctx context.Context, id string) (int64, error) {
	return u.Repository.DeleteThread(ctx, id)
}
