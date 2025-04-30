package bid

import (
	"context"
	"fmt"
	"renebizelli/go-leilao/Internals/entity/bid_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}
	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Bids not found, auctionId %v", auctionId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Bids not found, auctionId %v", auctionId))
	}

	var bids []BidEntityMongo

	if err := cursor.All(ctx, &bids); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode bids, auctionId %v", auctionId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Failed to decode bids, auctionId %v", auctionId))
	}

	var bidEntities []bid_entity.Bid

	for _, b := range bids {

		bidEntities = append(bidEntities, bid_entity.Bid{
			ID:        b.ID,
			AuctionID: b.AuctionID,
			UserID:    b.UserID,
			Amount:    b.Amount,
			Timestamp: time.Unix(b.Timestamp, 0),
		})
	}

	return bidEntities, nil
}

func (bd *BidRepository) FindWinning(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {

	fmt.Println("Bid repo 1:", auctionId)

	filter := bson.M{"auction_id": auctionId}

	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})

	var bidEntityMongo BidEntityMongo

	err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo)
	if err != nil {
		fmt.Printf("Failed to find winning bid, auctionId %v", auctionId)
		logger.Error(fmt.Sprintf("Failed to find winning bid, auctionId %v", auctionId), err)
		return nil, internal_error.NewInternalServerError(fmt.Sprintf("Failed to find winning bid, auctionId %v", auctionId))
	}

	fmt.Println("Bid repo 2:", bidEntityMongo.Amount)

	bid := &bid_entity.Bid{
		ID:        bidEntityMongo.ID,
		AuctionID: bidEntityMongo.AuctionID,
		UserID:    bidEntityMongo.UserID,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}

	return bid, nil
}
