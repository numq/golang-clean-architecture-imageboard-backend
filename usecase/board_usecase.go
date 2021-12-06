package usecase

import (
	"context"
	"grpcImageboard/domain"
)

type BoardUseCase interface {
	GetBoardsStream(ctx context.Context, changed func([]*domain.Board))
	GetBoards(ctx context.Context, skip int64, limit int64) ([]*domain.Board, error)
	GetBoardById(ctx context.Context, id string) (*domain.Board, error)
	CreateBoard(ctx context.Context, board *domain.Board) (string, error)
	DeleteBoard(ctx context.Context, id string) (int64, error)
}

type boardUseCase struct {
	Repository domain.BoardRepository
}

func NewBoardUseCase(repo domain.BoardRepository) BoardUseCase {
	return &boardUseCase{repo}
}

func (u *boardUseCase) GetBoardsStream(ctx context.Context, changed func([]*domain.Board)) {
	u.Repository.GetBoardsStream(ctx, changed)
}

func (u *boardUseCase) GetBoards(ctx context.Context, skip int64, limit int64) ([]*domain.Board, error) {
	return u.Repository.GetBoards(ctx, skip, limit)
}

func (u *boardUseCase) GetBoardById(ctx context.Context, id string) (*domain.Board, error) {
	return u.Repository.GetBoardById(ctx, id)
}

func (u *boardUseCase) CreateBoard(ctx context.Context, board *domain.Board) (string, error) {
	return u.Repository.CreateBoard(ctx, board)
}

func (u *boardUseCase) DeleteBoard(ctx context.Context, id string) (int64, error) {
	return u.Repository.DeleteBoard(ctx, id)
}
