package core

import (
	"context"
	"errors"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Customer struct {
	ID           int           `json:"id"`
	Limit        int           `json:"limite"`
	Balance      int           `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

type CustomerRepository interface {
	Tx(ctx context.Context, txFn func(context.Context) error) error
	Upsert(ctx context.Context, customer *Customer) error
	Find(ctx context.Context, customerID int, omitNested bool) (*Customer, error)
	Clear(ctx context.Context) error
}

func (c Customer) CheckBalance(value int) error {
	if c.Balance+c.Limit-value < 0 {
		return ErrInsufficientBalance
	}
	return nil
}

func Seed(ctx context.Context, customerRepo CustomerRepository) (func(context.Context) error, error) {
	customers := []Customer{
		{
			ID:           1,
			Balance:      0,
			Limit:        100000,
			Transactions: make([]Transaction, 0),
		},
		{
			ID:           2,
			Balance:      0,
			Limit:        80000,
			Transactions: make([]Transaction, 0),
		},
		{
			ID:           3,
			Balance:      0,
			Limit:        1000000,
			Transactions: make([]Transaction, 0),
		},
		{
			ID:           4,
			Balance:      0,
			Limit:        10000000,
			Transactions: make([]Transaction, 0),
		},
		{
			ID:           5,
			Balance:      0,
			Limit:        500000,
			Transactions: make([]Transaction, 0),
		},
	}

	for _, customer := range customers {
		if err := customerRepo.Upsert(ctx, &customer); err != nil {
			return nil, err
		}
	}

	return customerRepo.Clear, nil
}
