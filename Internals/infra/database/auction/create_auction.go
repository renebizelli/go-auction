package auction

import (
	"context"
	"os"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	ID          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) auction_entity.AuctionRepositoryInterface {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (repo *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auction := AuctionEntityMongo{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := repo.Collection.InsertOne(ctx, auction)
	if err != nil {
		logger.Error("Failed to create auction", err)
		return internal_error.NewInternalServerError("Failed to create auction")
	}

	go func() {

		<-time.After(getAuctionInterval())

		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		filter := bson.M{"_id": auctionEntity.ID}

		_, err := repo.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Failed to update auction status", err)
			return
		}
		logger.Info("Auction status updated to completed")

	}()

	return nil
}

func getAuctionInterval() time.Duration {

	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)

	if err != nil {
		logger.Error("Failed to parse auction interval", err)
		return time.Minute * 5
	}

	return duration
}
