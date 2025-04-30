package auction

import (
	"context"
	"errors"
	"fmt"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *AuctionRepository) FindAuctionByID(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {

	filter := bson.M{"_id": id}
	var auction AuctionEntityMongo
	err := repo.Collection.FindOne(ctx, filter).Decode(&auction)
	if err != nil {
		logger.Error(fmt.Sprintf("Auction %v not found", id), err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal_error.NewNotFoundError("Auction not found")
		}

		return nil, internal_error.NewInternalServerError("Failed to find auction")
	}

	auctionEntity := &auction_entity.Auction{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   time.Unix(auction.Timestamp, 0),
	}

	return auctionEntity, nil
}

func (repo *AuctionRepository) FindAuction(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {

	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{Pattern: productName, Options: "i"} // Case insensitive search
	}

	cursor, err := repo.Collection.Find(ctx, filter)

	if err != nil {
		logger.Error("Failed to find auctions", err)
		return nil, internal_error.NewInternalServerError("Failed to find auctions")
	}

	defer cursor.Close(ctx)

	var auctions []AuctionEntityMongo

	if err := cursor.All(ctx, &auctions); err != nil {
		logger.Error("Failed to decode auctions", err)
		return nil, internal_error.NewInternalServerError("Failed to find auctions")
	}

	var auctionEntities []auction_entity.Auction

	for _, auction := range auctions {
		a := auction_entity.Auction{
			ID:          auction.ID,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   auction.Condition,
			Status:      auction.Status,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		}

		auctionEntities = append(auctionEntities, a)
	}

	return auctionEntities, nil
}
