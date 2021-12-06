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

type postRepository struct {
	cfg config.Config
	db  db.Database
}

func CreatePostRepository(cfg config.Config, db db.Database) domain.PostRepository {
	return &postRepository{cfg, db}
}

func (r *postRepository) fetch(ctx context.Context, filter *bson.M, opts *options.FindOptions) ([]*domain.Post, error) {
	cursor, err := r.db.Posts.Find(ctx, filter, opts)
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
	data := make([]*domain.Post, 0)
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *postRepository) GetPosts(ctx context.Context, threadId string, skip int64, limit int64) ([]*domain.Post, error) {
	filter, opts := &bson.M{"thread_id": bson.M{"$eq": threadId}}, options.Find().SetSkip(skip).SetLimit(limit)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) (string, error) {
	post.Id = primitive.NewObjectID().Hex()
	_, err := r.db.Posts.InsertOne(ctx, post)
	if err != nil {
		return "", err
	}
	return post.Id, nil
}

func (r *postRepository) DeletePost(ctx context.Context, id string) (int64, error) {
	result, err := r.db.Posts.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
