package mongodb

import (
	"context"
	"errors"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type customerRepository struct {
	db *mongo.Collection
}

var _ core.CustomerRepository = (*customerRepository)(nil)

func NewCustomerRepository(db *mongo.Client) *customerRepository {
	return &customerRepository{
		db: db.Database("rinha").Collection("customers"),
	}
}

func (r *customerRepository) Tx(ctx context.Context, fn func(context.Context) error) (err error) {
	return Tx(ctx, r.db, fn)
}

func (r *customerRepository) Find(ctx context.Context, customerID int, omitNested bool) (*core.Customer, error) {
	var customer core.Customer

	opts := options.FindOne()
	if omitNested {
		opts.SetProjection(bson.D{{Key: "transactions", Value: 0}})
	}

	if err := r.db.FindOne(ctx, bson.D{{Key: "id", Value: customerID}}, opts).Decode(&customer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, core.ErrNotFound
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Upsert(ctx context.Context, customer *core.Customer) error {
	if _, err := r.Find(ctx, customer.ID, true); err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			return err
		}

		if _, err := r.db.InsertOne(ctx, customer); err != nil {
			return err
		}
	}

	return nil
}

func (r *customerRepository) Clear(ctx context.Context) error {
	if _, err := r.db.DeleteMany(ctx, bson.D{}); err != nil {
		return err
	}
	return nil
}
