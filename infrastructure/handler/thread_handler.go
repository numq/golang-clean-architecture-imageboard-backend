package handler

import (
	"context"
	"fmt"
	"grpcImageboard/mapping"
	"grpcImageboard/proto"
	"grpcImageboard/usecase"
	"io"
	"log"
)

type ThreadHandler interface {
	GetThreads(req *proto.GetThreadsReq, stream proto.ThreadService_GetThreadsServer) error
	GetHotThreads(_ *proto.GetHotThreadsReq, stream proto.ThreadService_GetHotThreadsServer) error
	GetLatestThreads(_ *proto.GetLatestThreadsReq, stream proto.ThreadService_GetLatestThreadsServer) error
	GetThreadById(ctx context.Context, req *proto.GetThreadByIdReq) (*proto.GetThreadByIdRes, error)
	CreateThread(ctx context.Context, req *proto.CreateThreadReq) (*proto.CreateThreadRes, error)
	UpdateThread(ctx context.Context, req *proto.UpdateThreadReq) (*proto.UpdateThreadRes, error)
	DeleteThread(ctx context.Context, req *proto.DeleteThreadReq) (*proto.DeleteThreadRes, error)
}

type threadHandler struct {
	useCase usecase.ThreadUseCase
}

func NewThreadHandler(uc usecase.ThreadUseCase) ThreadHandler {
	return &threadHandler{uc}
}

func (h *threadHandler) GetThreads(req *proto.GetThreadsReq, stream proto.ThreadService_GetThreadsServer) error {
	ctx := context.Background()
	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		return nil
	}
	data, err := h.useCase.GetThreads(ctx, req.GetBoardId(), req.GetSkip(), req.GetLimit())
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			return nil
		}
		return err
	}
	if data != nil {
		for _, thread := range data {
			if err := stream.Send(&proto.GetThreadsRes{Thread: mapping.ThreadToProto(thread)}); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	return nil
}

func (h *threadHandler) GetHotThreads(_ *proto.GetHotThreadsReq, stream proto.ThreadService_GetHotThreadsServer) error {
	ctx := context.Background()
	data, err := h.useCase.GetHotThreads(ctx)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			return nil
		}
	}
	if data != nil {
		for _, board := range data {
			if err := stream.Send(&proto.GetHotThreadsRes{Thread: mapping.ThreadToProto(board)}); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (h *threadHandler) GetLatestThreads(_ *proto.GetLatestThreadsReq, stream proto.ThreadService_GetLatestThreadsServer) error {
	ctx := context.Background()
	data, err := h.useCase.GetLatestThreads(ctx)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			return nil
		}
	}
	if data != nil {
		for _, board := range data {
			if err := stream.Send(&proto.GetLatestThreadsRes{Thread: mapping.ThreadToProto(board)}); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

func (h *threadHandler) GetThreadById(ctx context.Context, req *proto.GetThreadByIdReq) (*proto.GetThreadByIdRes, error) {
	thread, err := h.useCase.GetThreadById(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return &proto.GetThreadByIdRes{}, nil
	}
	return &proto.GetThreadByIdRes{Thread: mapping.ThreadToProto(thread)}, nil
}

func (h *threadHandler) CreateThread(ctx context.Context, req *proto.CreateThreadReq) (*proto.CreateThreadRes, error) {
	id, err := h.useCase.CreateThread(ctx, mapping.ThreadToModel(req.GetThread()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.CreateThreadRes{Id: id}, nil
}

func (h *threadHandler) UpdateThread(ctx context.Context, req *proto.UpdateThreadReq) (*proto.UpdateThreadRes, error) {
	thread, err := h.useCase.UpdateThread(ctx, mapping.ThreadToModel(req.GetThread()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.UpdateThreadRes{Thread: mapping.ThreadToProto(thread)}, nil
}

func (h *threadHandler) DeleteThread(ctx context.Context, req *proto.DeleteThreadReq) (*proto.DeleteThreadRes, error) {
	count, err := h.useCase.DeleteThread(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.DeleteThreadRes{Count: count}, nil
}
