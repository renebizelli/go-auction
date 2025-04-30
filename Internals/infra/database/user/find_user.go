package user

import (
	"context"
	"errors"
	"fmt"
	"renebizelli/go-leilao/Internals/entity/user_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *UserRepository) FindUserByID(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {

	filter := bson.M{"_id": userId}

	var userEntityMongo UserEntityMongo

	err := r.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)

	if err != nil {

		logger.Error(fmt.Sprintf("Error finding user by ID: %v", userId), err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal_error.NewNotFoundError(fmt.Sprintf("User %v not found", userId))
		}

		return nil, internal_error.NewInternalServerError(err.Error())
	}

	user := &user_entity.User{
		ID:   userEntityMongo.ID,
		Name: userEntityMongo.Name,
	}

	return user, nil
}
