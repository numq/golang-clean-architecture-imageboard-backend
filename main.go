package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpcImageboard/config"
	"grpcImageboard/data/db"
	"grpcImageboard/data/repository"
	"grpcImageboard/infrastructure"
	"grpcImageboard/infrastructure/handler"
	"grpcImageboard/proto"
	"grpcImageboard/usecase"
)

func main() {

	cfg := config.NewConfig()

	client, connect, disconnect := db.Open(fmt.Sprintf("%s/%s", cfg.Database.DbUri, cfg.Database.DbName))
	defer db.Close(client, connect, disconnect)

	appDatabase := db.NewDatabase(*cfg, client.Database(cfg.Database.DbName))

	boardRepo, threadRepo, postRepo := repository.CreateBoardRepository(*cfg, appDatabase), repository.CreateThreadRepository(*cfg, appDatabase), repository.CreatePostRepository(*cfg, appDatabase)

	boardUseCase, threadUseCase, postUseCase := usecase.NewBoardUseCase(boardRepo), usecase.NewThreadUseCase(threadRepo), usecase.NewPostUseCase(postRepo)

	boardHandler := handler.NewBoardHandler(boardUseCase)
	threadHandler := handler.NewThreadHandler(threadUseCase)
	postHandler := handler.NewPostHandler(postUseCase)

	infrastructure.Init(cfg.Server.ServerHost, cfg.Server.ServerPort, func(server *grpc.Server) {
		proto.RegisterBoardServiceServer(server, boardHandler)
		proto.RegisterThreadServiceServer(server, threadHandler)
		proto.RegisterPostServiceServer(server, postHandler)
	})

}
