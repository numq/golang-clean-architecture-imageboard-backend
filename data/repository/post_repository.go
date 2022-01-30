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
	"sort"
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
	filter, opts := &bson.M{"thread_id": bson.M{"$eq": threadId}}, options.Find().SetSort(bson.M{"created_at": -1}).SetSkip(skip).SetLimit(limit)
	data, err := r.fetch(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *postRepository) GetPostById(ctx context.Context, id string) (*domain.Post, error) {
	data := new(domain.Post)
	if err := r.db.Posts.FindOne(ctx, bson.M{"_id": id}).Decode(&data); err != nil {
		return nil, err
	}
	if data.Id != "" {
		return data, nil
	}
	return nil, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) (string, error) {
	session, err := r.db.Boards.Database().Client().StartSession()
	if err != nil {
		return "", err
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		post.Id = primitive.NewObjectID().Hex()
		ids := post.QuoteIds
		sort.Strings(ids)
		post.QuoteIds = ids
		if _, err := r.db.Posts.InsertOne(ctx, post); err != nil {
			return err
		}
		count, err := r.db.Posts.CountDocuments(ctx, bson.M{"thread_id": post.ThreadId})
		if err != nil {
			return err
		}
		if _, err = r.db.Threads.UpdateOne(ctx, bson.M{"_id": post.ThreadId}, bson.M{"$set": bson.M{"post_count": count}}); err != nil {
			return err
		}
		return nil
	}); err != nil {
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
