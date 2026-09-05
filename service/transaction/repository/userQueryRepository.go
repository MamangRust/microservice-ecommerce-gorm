package repository

import (
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
)

type userQueryRepository struct {
	client pbuser.UserQueryServiceClient
}

func NewUserQueryRepository(client pbuser.UserQueryServiceClient) *userQueryRepository {
	return &userQueryRepository{
		client: client,
	}
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*dto.GetUserByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbuser.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
	}

	return &dto.GetUserByIDRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}
