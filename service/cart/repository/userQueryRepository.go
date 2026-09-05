package repository

import (
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
)

type userQueryRepository struct {
	client pbuser.UserQueryServiceClient
}

func NewUserQueryRepository(client pbuser.UserQueryServiceClient) UserQueryRepository {
	return &userQueryRepository{client: client}
}

func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*UserResult, error) {
	res, err := r.client.FindById(ctx, &pbuser.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &UserResult{
		UserID: res.Data.Id,
	}, nil
}
