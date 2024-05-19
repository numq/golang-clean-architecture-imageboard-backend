package mapper

import (
	"grpcImageboard/domain"
	"grpcImageboard/proto"
)

func PostToProto(post *domain.Post) *proto.Post {
	return &proto.Post{
		Id:          post.Id,
		ThreadId:    post.ThreadId,
		Description: post.Description,
		QuoteIds:    post.QuoteIds,
		Text:        post.Text,
		ImageUrl:    post.ImageUrl,
		VideoUrl:    post.VideoUrl,
		CreatedAt:   post.CreatedAt,
	}
}

func PostToModel(post *proto.Post) *domain.Post {
	return &domain.Post{
		Id:          post.GetId(),
		ThreadId:    post.GetThreadId(),
		Description: post.GetDescription(),
		QuoteIds:    post.GetQuoteIds(),
		Text:        post.GetText(),
		ImageUrl:    post.GetImageUrl(),
		VideoUrl:    post.GetVideoUrl(),
		CreatedAt:   post.GetCreatedAt(),
	}
}
