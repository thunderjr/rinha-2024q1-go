package core

import (
	"context"
	"time"
)

type Transaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

type TransactionRepository interface {
	Tx(ctx context.Context, txFn func(context.Context) error) error
	Create(ctx context.Context, customerID int, transaction *Transaction) error
}
