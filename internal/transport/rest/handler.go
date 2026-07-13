package rest

import (
	"desafio-clean-architecture/internal/domain"
	"desafio-clean-architecture/internal/usecase"
	"encoding/json"
	"net/http"
)

type Handler struct {
	listUseCase   *usecase.ListOrdersUseCase
	createUseCase *usecase.CreateOrderUseCase
}

func NewHandler(listUseCase *usecase.ListOrdersUseCase, createUseCase *usecase.CreateOrderUseCase) *Handler {
	return &Handler{listUseCase: listUseCase, createUseCase: createUseCase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.list(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	orders, err := h.listUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"orders": orders})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var input domain.Order
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.createUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
