package repository

import (
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-merchant/dto"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
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
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &dto.GetUserByIDRow{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  "", // Not provided by gRPC, but required by struct
	}, nil
}
