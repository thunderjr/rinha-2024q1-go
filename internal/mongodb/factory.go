package mongodb

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context) (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	return mongo.Connect(ctx, options.Client().ApplyURI(uri).SetDirect(true))
}
