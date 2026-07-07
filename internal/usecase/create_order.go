package usecase

import "desafio-clean-architecture/internal/domain"

type CreateOrderUseCase struct {
	repo domain.OrderRepository
}

func NewCreateOrderUseCase(repo domain.OrderRepository) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: repo}
}

func (u *CreateOrderUseCase) Execute(input domain.Order) (domain.Order, error) {
	if input.CustomerName == "" {
		input.CustomerName = "Anonymous"
	}
	if input.Status == "" {
		input.Status = "created"
	}
	return u.repo.Create(input)
}
