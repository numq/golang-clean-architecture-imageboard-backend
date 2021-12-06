package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpcImageboard/config"
	"log"
	"time"
)

type Database struct {
	Boards  *mongo.Collection
	Threads *mongo.Collection
	Posts   *mongo.Collection
}

func NewDatabase(cfg config.Config, db *mongo.Database) Database {
	return Database{
		Boards:  db.Collection(cfg.Database.CollBoards),
		Threads: db.Collection(cfg.Database.CollThreads),
		Posts:   db.Collection(cfg.Database.CollPosts),
	}
}

func Open(uri string) (*mongo.Client, context.Context, context.CancelFunc) {
	return connect(uri)
}

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	disconnect(client, ctx, cancel)
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetDirect(true).SetReplicaSet("rs0"))
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Connected to DB")
	}
	return client, ctx, cancel
}

func disconnect(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println(err)
		} else {
			log.Println("Disconnected from DB")
		}
	}()
}
