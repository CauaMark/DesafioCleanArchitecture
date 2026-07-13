package grpcserver

import (
	"context"
	orderspb "desafio-clean-architecture/internal/proto/orderspb"
	"desafio-clean-architecture/internal/usecase"
)

type Server struct {
	orderspb.UnimplementedOrderServiceServer
	listUseCase *usecase.ListOrdersUseCase
}

func NewServer(listUseCase *usecase.ListOrdersUseCase) *Server {
	return &Server{listUseCase: listUseCase}
}

func (s *Server) ListOrders(ctx context.Context, req *orderspb.ListOrdersRequest) (*orderspb.ListOrdersResponse, error) {
	_ = ctx
	_ = req
	orders, err := s.listUseCase.Execute()
	if err != nil {
		return nil, err
	}

	resp := &orderspb.ListOrdersResponse{}
	for _, order := range orders {
		resp.Orders = append(resp.Orders, &orderspb.Order{
			Id:           order.ID,
			CustomerName: order.CustomerName,
			Total:        order.Total,
			Status:       order.Status,
		})
	}
	return resp, nil
}
