package mapper

import (
	"grpcImageboard/domain"
	"grpcImageboard/proto"
)

func ThreadToProto(thread *domain.Thread) *proto.Thread {
	return &proto.Thread{
		Id:        thread.Id,
		BoardId:   thread.BoardId,
		PostCount: thread.PostCount,
		Title:     thread.Title,
		CreatedAt: thread.CreatedAt,
		BumpedAt:  thread.BumpedAt,
	}
}

func ThreadToModel(thread *proto.Thread) *domain.Thread {
	return &domain.Thread{
		Id:        thread.GetId(),
		BoardId:   thread.GetBoardId(),
		PostCount: thread.GetPostCount(),
		Title:     thread.GetTitle(),
		CreatedAt: thread.GetCreatedAt(),
		BumpedAt:  thread.GetBumpedAt(),
	}
}
