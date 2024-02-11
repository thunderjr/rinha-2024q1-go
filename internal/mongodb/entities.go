package mongodb

import (
	"time"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
)

type TransactionEntity struct {
	Value       int
	Type        string
	Description string
	CreatedAt   time.Time
}

func (t TransactionEntity) ToModel() core.Transaction {
	return core.Transaction{
		Value:       t.Value,
		Type:        t.Type,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}

type CustomerEntity struct {
	ID           int
	Limit        int
	Balance      int
	Transactions []TransactionEntity
}

func (c CustomerEntity) ToModel() core.Customer {
	txs := make([]core.Transaction, 0, len(c.Transactions))
	for _, tx := range c.Transactions {
		txs = append(txs, tx.ToModel())
	}

	return core.Customer{
		Transactions: txs,
		ID:           c.ID,
		Limit:        c.Limit,
		Balance:      c.Balance,
	}
}
