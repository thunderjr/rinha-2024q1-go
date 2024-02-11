package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
)

type extractResponse struct {
	Saldo struct {
		Total int    `json:"total"`
		Limit int    `json:"limite"`
		Date  string `json:"data_extrato"`
	} `json:"saldo"`
	LastTransactions []core.Transaction `json:"ultimas_transacoes"`
}

func (s *server) GetExtract(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	customerIDPath := req.PathValue("id")
	customerID, err := strconv.Atoi(customerIDPath)
	if err != nil {
		NewError(w, err).Status(http.StatusBadRequest).Send()
		return
	}

	customer, err := s.deps.CustomerRepo.Find(ctx, customerID, false)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			NewError(w, err).Status(http.StatusNotFound).Send()
			return
		}

		NewError(w, err).Status(http.StatusInternalServerError).Send()
		return
	}

	response := extractResponse{
		LastTransactions: customer.Transactions,
	}

	response.Saldo.Limit = customer.Limit
	response.Saldo.Total = customer.Balance
	response.Saldo.Date = time.Now().UTC().Format(time.RFC3339Nano)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		NewError(w, err).Status(http.StatusInternalServerError).Send()
		return
	}
}
