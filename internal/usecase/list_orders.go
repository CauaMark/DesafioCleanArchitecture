package usecase

import "desafio-clean-architecture/internal/domain"

type ListOrdersUseCase struct {
	repo domain.OrderRepository
}

func NewListOrdersUseCase(repo domain.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{repo: repo}
}

func (u *ListOrdersUseCase) Execute() ([]domain.Order, error) {
	return u.repo.List()
}
