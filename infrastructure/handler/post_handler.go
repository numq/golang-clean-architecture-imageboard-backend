package handler

import (
	"context"
	"grpcImageboard/mapper"
	"grpcImageboard/proto"
	"grpcImageboard/usecase"
	"log"
)

type PostHandler interface {
	GetPosts(req *proto.GetPostsReq, stream proto.PostService_GetPostsServer) error
	GetPostById(ctx context.Context, req *proto.GetPostByIdReq) (*proto.GetPostByIdRes, error)
	CreatePost(ctx context.Context, req *proto.CreatePostReq) (*proto.CreatePostRes, error)
	DeletePost(ctx context.Context, req *proto.DeletePostReq) (*proto.DeletePostRes, error)
}

type postHandler struct {
	proto.UnimplementedPostServiceServer
	useCase usecase.PostUseCase
}

func NewPostHandler(uc usecase.PostUseCase) proto.PostServiceServer {
	return &postHandler{useCase: uc}
}

func (h *postHandler) GetPosts(req *proto.GetPostsReq, stream proto.PostService_GetPostsServer) error {
	ctx := context.Background()
	data, err := h.useCase.GetPosts(ctx, req.GetThreadId(), req.GetSkip(), req.GetLimit())
	if err != nil {
		log.Println(err)
	}
	if data != nil {
		for _, post := range data {
			if err := stream.Send(&proto.GetPostsRes{Post: mapper.PostToProto(post)}); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (h *postHandler) GetPostById(ctx context.Context, req *proto.GetPostByIdReq) (*proto.GetPostByIdRes, error) {
	post, err := h.useCase.GetPostById(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return &proto.GetPostByIdRes{}, nil
	}
	return &proto.GetPostByIdRes{Post: mapper.PostToProto(post)}, nil
}

func (h *postHandler) CreatePost(ctx context.Context, req *proto.CreatePostReq) (*proto.CreatePostRes, error) {
	id, err := h.useCase.CreatePost(ctx, mapper.PostToModel(req.GetPost()))
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
