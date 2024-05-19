package mapper

import (
	"grpcImageboard/domain"
	"grpcImageboard/proto"
)

func BoardToProto(board *domain.Board) *proto.Board {
	return &proto.Board{
		Id:          board.Id,
		Title:       board.Title,
		Description: board.Description,
		ImageUrl:    board.ImageUrl,
		IsAdult:     board.IsAdult,
	}
}

func BoardToModel(board *proto.Board) *domain.Board {
	return &domain.Board{
		Id:          board.GetId(),
		Title:       board.GetTitle(),
		Description: board.GetDescription(),
		ImageUrl:    board.GetImageUrl(),
		IsAdult:     board.GetIsAdult(),
	}
}
