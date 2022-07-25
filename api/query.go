package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IQuery - BaseQuery implementations
type IQuery interface {
	find(
		ctx context.Context,
		collectionName string,
		filter bson.M,
		opts ...*options.FindOptions,
	) []bson.M
}

// BaseQuery - hold db
type BaseQuery struct {
	db *mongo.Database
}

// NewBaseQuery - instantiate a new BaseQuery
func NewBaseQuery(db *mongo.Database) *BaseQuery {
	return &BaseQuery{db: db}
}

// Find - find all documents matching filter
func (b *BaseQuery) find(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []bson.M {
	collection := b.db.Collection(collectionName)

	cursor, _ := collection.Find(ctx, filter)

	var decodedResult []bson.M
	cursor.All(context.TODO(), &decodedResult)

	return decodedResult
}
