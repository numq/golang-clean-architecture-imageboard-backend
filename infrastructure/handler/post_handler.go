package handler

import (
	"context"
	"grpcImageboard/mapping"
	"grpcImageboard/proto"
	"grpcImageboard/usecase"
	"log"
)

type PostHandler interface {
	GetPosts(req *proto.GetPostsReq, stream proto.PostService_GetPostsServer) error
	CreatePost(ctx context.Context, req *proto.CreatePostReq) (*proto.CreatePostRes, error)
	DeletePost(ctx context.Context, req *proto.DeletePostReq) (*proto.DeletePostRes, error)
}

type postHandler struct {
	useCase usecase.PostUseCase
}

func NewPostHandler(uc usecase.PostUseCase) PostHandler {
	return &postHandler{uc}
}

func (h *postHandler) GetPosts(req *proto.GetPostsReq, stream proto.PostService_GetPostsServer) error {
	ctx := context.Background()
	data, err := h.useCase.GetPosts(ctx, req.GetThreadId(), req.GetSkip(), req.GetLimit())
	if err != nil {
		log.Println(err)
	}
	if data != nil {
		for _, post := range data {
			if err := stream.Send(&proto.GetPostsRes{Post: mapping.PostToProto(post)}); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (h *postHandler) CreatePost(ctx context.Context, req *proto.CreatePostReq) (*proto.CreatePostRes, error) {
	id, err := h.useCase.CreatePost(ctx, mapping.PostToModel(req.GetPost()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.CreatePostRes{Id: id}, nil
}

func (h *postHandler) DeletePost(ctx context.Context, req *proto.DeletePostReq) (*proto.DeletePostRes, error) {
	count, err := h.useCase.DeletePost(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &proto.DeletePostRes{Count: count}, nil
}
