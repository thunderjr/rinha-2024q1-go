package httpx

import (
	"fmt"
	"log"
	"net/http"

	"github.com/thunderjr/rinha-2024q1-go/internal/core"
)

type Deps struct {
	TransactionRepo core.TransactionRepository
	CustomerRepo    core.CustomerRepository
}

type server struct {
	deps Deps
}

func New(port string, deps Deps) error {
	server := server{deps}
	m := http.NewServeMux()

	m.HandleFunc("POST /clientes/{id}/transacoes", server.CreateTransaciton)
	m.HandleFunc("GET /clientes/{id}/extrato", server.GetExtract)

	log.Printf("Server running on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), m)
}
