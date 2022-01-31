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

type threadRepository struct {
	cfg config.Config
	db  db.Database
}

func CreateThreadRepository(cfg config.Config, db db.Database) domain.ThreadRepository {
	return &threadRepository{cfg, db}
}

func (r *threadRepository) fetch(ctx context.Context, filter *bson.M, opts *options.FindOptions) ([]*domain.Thread, error) {
	cursor, err := r.db.Threads.Find(ctx, filter, opts)
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
	data := make([]*domain.Thread, 0)
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *threadRepository) GetThreads(ctx context.Context, boardId string, skip int64, limit int64) ([]*domain.Thread, error) {
	filter, opts := &bson.M{"board_id": bson.M{"$eq": boardId}}, options.Find().SetSort(bson.M{"bumped_at": -1}).SetSkip(skip).SetLimit(limit)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *threadRepository) GetHotThreads(ctx context.Context) ([]*domain.Thread, error) {
	filter, opts := &bson.M{}, options.Find().SetSort(bson.M{"post_count": -1}).SetLimit(10)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *threadRepository) GetLatestThreads(ctx context.Context) ([]*domain.Thread, error) {
	filter, opts := &bson.M{}, options.Find().SetSort(bson.M{"bumped_at": -1}).SetLimit(10)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *threadRepository) GetThreadById(ctx context.Context, id string) (*domain.Thread, error) {
	data := new(domain.Thread)
	if err := r.db.Threads.FindOne(ctx, bson.M{"_id": id}).Decode(&data); err != nil {
		return nil, err
	}
	if data.Id != "" {
		return data, nil
	}
	return nil, nil
}

func (r *threadRepository) CreateThread(ctx context.Context, thread *domain.Thread) (string, error) {
	thread.Id = primitive.NewObjectID().Hex()
	_, err := r.db.Threads.InsertOne(ctx, thread)
	if err != nil {
		return "", err
	}
	return thread.Id, nil
}

func (r *threadRepository) UpdateThread(ctx context.Context, thread *domain.Thread) (*domain.Thread, error) {
	result := r.db.Threads.FindOneAndUpdate(ctx, bson.M{"_id": thread.Id}, thread)
	if result != nil {
		var data *domain.Thread
		if err := result.Decode(&data); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func (r *threadRepository) DeleteThread(ctx context.Context, id string) (int64, error) {
	session, err := r.db.Threads.Database().Client().StartSession()
	if err != nil {
		return 0, err
	}
	var count int64 = 0
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		result, err := r.db.Threads.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
		count = result.DeletedCount
		if _, err = r.db.Posts.DeleteMany(ctx, bson.M{"thread_id": id}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}
