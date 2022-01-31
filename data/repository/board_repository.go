package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpcImageboard/config"
	"grpcImageboard/data/db"
	"grpcImageboard/domain"
)

type boardRepository struct {
	cfg config.Config
	db  db.Database
}

func CreateBoardRepository(cfg config.Config, db db.Database) domain.BoardRepository {
	return &boardRepository{cfg, db}
}

func (r *boardRepository) GetBoardsStream(ctx context.Context, changed func([]*domain.Board)) {
	stream, err := r.db.Boards.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		fmt.Println(err)
	}
	defer func(stream *mongo.ChangeStream, ctx context.Context) {
		err := stream.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(stream, ctx)
	for {
		select {
		case <-ctx.Done():
			err := stream.Close(ctx)
			if err != nil {
				fmt.Println(err)
			}
			return
		default:
			var data []*domain.Board
			ok := stream.Next(ctx)
			if !ok {
				err := stream.Err()
				if err != nil {
					println(err)
					return
				}
			}
			if ok {
				err := stream.Decode(data)
				if err != nil {
					println(err)
					return
				}
				changed(data)
			}
		}
	}
}

func (r *boardRepository) fetch(ctx context.Context, filter *bson.M, opts *options.FindOptions) ([]*domain.Board, error) {
	cursor, err := r.db.Boards.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
	}(cursor, ctx)
	data := make([]*domain.Board, 0)
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *boardRepository) GetBoards(ctx context.Context, skip int64, limit int64) ([]*domain.Board, error) {
	filter, opts := &bson.M{}, options.Find().SetSort(bson.M{"title": 1}).SetSkip(skip).SetLimit(limit)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *boardRepository) GetBoardById(ctx context.Context, id string) (*domain.Board, error) {
	data := new(domain.Board)
	if err := r.db.Boards.FindOne(ctx, bson.M{"_id": bson.M{"$eq": id}}).Decode(&data); err != nil {
		return nil, err
	}
	if data.Id != "" {
		return data, nil
	}
	return nil, nil
}

func (r *boardRepository) CreateBoard(ctx context.Context, board *domain.Board) (string, error) {
	board.Id = primitive.NewObjectID().Hex()
	_, err := r.db.Boards.InsertOne(ctx, board)
	if err != nil {
		return "", err
	}
	return board.Id, nil
}

func (r *boardRepository) DeleteBoard(ctx context.Context, id string) (int64, error) {
	session, err := r.db.Boards.Database().Client().StartSession()
	if err != nil {
		return 0, err
	}
	var count int64 = 0
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		result, err := r.db.Boards.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
		count = result.DeletedCount
		thread := domain.Thread{}
		if err := r.db.Threads.FindOne(ctx, bson.M{"board_id": id}).Decode(&thread); err != nil {
			return err
		}
		if _, err = r.db.Threads.DeleteMany(ctx, bson.D{{"board_id", id}}); err != nil {
			return err
		}
		if _, err = r.db.Posts.DeleteMany(ctx, bson.D{{"thread_id", thread.Id}}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}
