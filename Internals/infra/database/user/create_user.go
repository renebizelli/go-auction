package user

import (
	"context"
	"renebizelli/go-leilao/Internals/entity/user_entity"
	"renebizelli/go-leilao/Internals/internal_error"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) user_entity.UserRepositoryInterface {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *user_entity.User) *internal_error.InternalError {

	userEntityMongo := &UserEntityMongo{
		ID:   user.ID,
		Name: user.Name,
	}

	_, err := u.Collection.InsertOne(ctx, userEntityMongo)
	if err != nil {
		return internal_error.NewInternalServerError(err.Error())
	}

	return nil
}
