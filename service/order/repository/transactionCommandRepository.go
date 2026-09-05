package repository

import (
	"context"
	"fmt"

		pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandRepository struct {
	client pbtransaction.TransactionCommandServiceClient
}

func NewTransactionCommandRepository(client pbtransaction.TransactionCommandServiceClient) TransactionCommandRepository {
	return &transactionCommandRepository{client: client}
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	_, err := r.client.DeleteTransactionByOrderPermanent(ctx, &pbtransaction.FindByIdTransactionRequest{
		Id: int32(order_id),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("transaction command client is not initialized")
	}
	res, err := r.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}
