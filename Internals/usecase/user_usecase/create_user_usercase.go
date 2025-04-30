package user_usecase

import (
	"context"
	"renebizelli/go-leilao/Internals/entity/user_entity"
	"renebizelli/go-leilao/Internals/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserInputDTO struct {
	Name string `json:"name"`
}

type UserOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserByID(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
	CreateUser(ctx context.Context, userInput UserInputDTO) *internal_error.InternalError
}

func NewUserUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		UserRepository: userRepository,
	}
}

func (a *UserUseCase) CreateUser(ctx context.Context, userInput UserInputDTO) *internal_error.InternalError {

	userEntity := user_entity.CreateUser(userInput.Name)

	if e := a.UserRepository.CreateUser(ctx, userEntity); e != nil {
		return e
	}

	return nil
}
