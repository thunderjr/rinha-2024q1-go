package mongodb

import (
	"context"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionRepository struct {
	db *mongo.Collection
}

var _ core.TransactionRepository = (*transactionRepository)(nil)

func NewTransactionRepository(db *mongo.Client) *transactionRepository {
	return &transactionRepository{
		db: db.Database("rinha").Collection("customers"),
	}
}

func (r *transactionRepository) Tx(ctx context.Context, fn func(context.Context) error) (err error) {
	return Tx(ctx, r.db, fn)
}

func (r *transactionRepository) Create(ctx context.Context, customerID int, transaction *core.Transaction) error {
	customerFilter := bson.D{{Key: "id", Value: customerID}}

	customerPushTransaction := bson.D{{Key: "$push", Value: bson.D{{Key: "transactions", Value: transaction}}}}
	customerUpdateBalance := bson.D{{Key: "$inc", Value: bson.D{{Key: "balance", Value: transaction.Value}}}}

	ops := []interface{}{customerPushTransaction, customerUpdateBalance}
	for _, op := range ops {
		if _, err := r.db.UpdateOne(ctx, customerFilter, op); err != nil {
			return err
		}
	}

	return nil
}
