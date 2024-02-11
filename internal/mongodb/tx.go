package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func Tx(ctx context.Context, db *mongo.Collection, fn func(context.Context) error) (err error) {
	session, err := db.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if _, err = session.WithTransaction(ctx,
		func(ctx mongo.SessionContext) (interface{}, error) {
			return nil, fn(ctx)
		},
		options.Transaction().SetWriteConcern(writeconcern.Majority()),
	); err != nil {
		return err
	}
	return nil
}
