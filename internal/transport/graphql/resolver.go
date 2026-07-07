package graphql

import (
	"desafio-clean-architecture/internal/usecase"
	"encoding/json"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"net/http"
)

type Resolver struct {
	listUseCase *usecase.ListOrdersUseCase
}

func NewResolver(listUseCase *usecase.ListOrdersUseCase) *Resolver {
	return &Resolver{listUseCase: listUseCase}
}

type Order struct {
	ID           graphql.ID `json:"id"`
	CustomerName string     `json:"customerName"`
	Total        float64    `json:"total"`
	Status       string     `json:"status"`
}

func (r *Resolver) ListOrders() []Order {
	orders, err := r.listUseCase.Execute()
	if err != nil {
		return nil
	}

	result := make([]Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, Order{
			ID:           graphql.ID(fmt.Sprintf("%d", order.ID)),
			CustomerName: order.CustomerName,
			Total:        order.Total,
			Status:       order.Status,
		})
	}
	return result
}

func NewHandler(listUseCase *usecase.ListOrdersUseCase) http.Handler {
	schemaString := `
schema {
  query: Query
}

type Query {
  listOrders: [Order!]!
}

type Order {
  id: ID!
  customerName: String!
  total: Float!
  status: String!
}
`
	parsed := graphql.MustParseSchema(schemaString, NewResolver(listUseCase))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		if r.Method == http.MethodGet {
			request.Query = r.URL.Query().Get("query")
		} else {
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
		}

		if request.Query == "" {
			http.Error(w, "missing query", http.StatusBadRequest)
			return
		}

		response := parsed.Exec(r.Context(), request.Query, "", request.Variables)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})
}
