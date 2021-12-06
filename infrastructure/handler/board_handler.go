package handler

import (
	"context"
	"grpcImageboard/mapping"
	"grpcImageboard/proto"
	"grpcImageboard/usecase"
	"io"
	"log"
)

type BoardHandler interface {
	GetBoards(req *proto.GetBoardsReq, stream proto.BoardService_GetBoardsServer) error
	GetBoardById(ctx context.Context, req *proto.GetBoardByIdReq) (*proto.GetBoardByIdRes, error)
	CreateBoard(ctx context.Context, req *proto.CreateBoardReq) (*proto.CreateBoardRes, error)
	DeleteBoard(ctx context.Context, req *proto.DeleteBoardReq) (*proto.DeleteBoardRes, error)
}

type boardHandler struct {
	useCase usecase.BoardUseCase
}

func NewBoardHandler(uc usecase.BoardUseCase) BoardHandler {
	return &boardHandler{uc}
}

func (h *boardHandler) GetBoards(req *proto.GetBoardsReq, stream proto.BoardService_GetBoardsServer) error {
	ctx := context.Background()
	data, err := h.useCase.GetBoards(ctx, req.GetSkip(), req.GetLimit())
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			return nil
		}
	}
	if data != nil {
		for _, board := range data {
			if err := stream.Send(&proto.GetBoardsRes{Board: mapping.BoardToProto(board)}); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (h *boardHandler) GetBoardById(ctx context.Context, req *proto.GetBoardByIdReq) (*proto.GetBoardByIdRes, error) {
	board, err := h.useCase.GetBoardById(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return &proto.GetBoardByIdRes{}, nil
	}
	return &proto.GetBoardByIdRes{Board: mapping.BoardToProto(board)}, nil
}

func (h *boardHandler) CreateBoard(ctx context.Context, req *proto.CreateBoardReq) (*proto.CreateBoardRes, error) {
	id, err := h.useCase.CreateBoard(ctx, mapping.BoardToModel(req.GetBoard()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.CreateBoardRes{Id: id}, nil
}

func (h *boardHandler) DeleteBoard(ctx context.Context, req *proto.DeleteBoardReq) (*proto.DeleteBoardRes, error) {
	count, err := h.useCase.DeleteBoard(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.DeleteBoardRes{Count: count}, nil
}
