package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
)

var (
	ErrInvalidValue       = errors.New("invalid value")
	ErrInvalidType        = errors.New("invalid type")
	ErrInvalidDescription = errors.New("invalid description")
)

func (s *server) CreateTransaciton(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	customerIDPath := req.PathValue("id")
	customerID, err := strconv.Atoi(customerIDPath)
	if err != nil {
		NewError(w, err).Status(http.StatusBadRequest).Send()
		return
	}

	transaction := &core.Transaction{}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(req.Body).Decode(transaction); err != nil {
		NewError(w, err).Status(http.StatusBadRequest).Send()
		return
	}

	if transaction.Value <= 0 || transaction.Value == 0 {
		NewError(w, ErrInvalidValue).Status(http.StatusBadRequest).Send()
		return
	}

	if transaction.Description == "" || len(transaction.Description) > 10 {
		NewError(w, ErrInvalidDescription).Status(http.StatusBadRequest).Send()
		return
	}

	if transaction.Type != "c" && transaction.Type != "d" {
		NewError(w, ErrInvalidType).Status(http.StatusBadRequest).Send()
		return
	}

	var customer *core.Customer
	txFn := func(ctx context.Context) error {
		customer, err = s.deps.CustomerRepo.Find(ctx, customerID, true)
		if err != nil {
			return err
		}

		if transaction.Type == "d" {
			if err := customer.CheckBalance(transaction.Value); err != nil {
				return err
			}

			transaction.Value = -transaction.Value
		}

		transaction.CreatedAt = time.Now()
		if err := s.deps.TransactionRepo.Create(ctx, customerID, transaction); err != nil {
			return err
		}

		return nil
	}

	if err := s.deps.TransactionRepo.Tx(ctx, txFn); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			NewError(w, err).Status(http.StatusNotFound).Send()
			return
		}

		if errors.Is(err, core.ErrInsufficientBalance) {
			NewError(w, err).Status(http.StatusUnprocessableEntity).Send()
			return
		}

		NewError(w, err).Status(http.StatusInternalServerError).Send()
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		NewError(w, err).Status(http.StatusInternalServerError).Send()
		return
	}
}
