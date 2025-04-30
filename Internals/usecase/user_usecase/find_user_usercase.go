package user_usecase

import (
	"context"
	"renebizelli/go-leilao/Internals/internal_error"
)

func (u *UserUseCase) FindUserByID(
	ctx context.Context,
	id string) (*UserOutputDTO, *internal_error.InternalError) {

	user, err := u.UserRepository.FindUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	userOutput := &UserOutputDTO{
		ID:   user.ID,
		Name: user.Name,
	}

	return userOutput, nil
}
