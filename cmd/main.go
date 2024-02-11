package main

import (
	"context"
	"log"
	"os"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
	"github.com/thunderjr/rinha-2024q1-go/internal/httpx"
	"github.com/thunderjr/rinha-2024q1-go/internal/mongodb"
)

func main() {
	ctx := context.Background()

	db, err := mongodb.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect(ctx)

	customerRepo := mongodb.NewCustomerRepository(db)
	transactionRepo := mongodb.NewTransactionRepository(db)

	clearDb, err := core.Seed(ctx, customerRepo)
	if err != nil {
		log.Fatal(err)
	}
	defer clearDb(ctx)

	port := os.Getenv("PORT")
	log.Fatal(httpx.New(port, httpx.Deps{
		TransactionRepo: transactionRepo,
		CustomerRepo:    customerRepo,
	}))
}
