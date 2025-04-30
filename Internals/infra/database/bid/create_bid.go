package bid

import (
	"context"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"renebizelli/go-leilao/Internals/entity/bid_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	ID        string  `bson:"_id"`
	AuctionID string  `bson:"auction_id"`
	UserID    string  `bson:"user_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection  *mongo.Collection
	AuctionRepo auction_entity.AuctionRepositoryInterface
}

func NewBidRepository(database *mongo.Database, auctionRepository auction_entity.AuctionRepositoryInterface) *BidRepository {
	return &BidRepository{
		Collection:  database.Collection("bids"),
		AuctionRepo: auctionRepository,
	}
}

func (bd *BidRepository) CreateBid(ctx context.Context, bids []bid_entity.Bid) *internal_error.InternalError {

	wg := sync.WaitGroup{}

	for _, bid := range bids {
		wg.Add(1)

		go func(bidValue bid_entity.Bid) {

			defer wg.Done()

			auction, err := bd.AuctionRepo.FindAuctionByID(ctx, bidValue.AuctionID)

			if err != nil {
				logger.Error("Failed to find auction", err)
				return
			}

			if auction.Status != auction_entity.Active {
				return
			}

			bidEntityMongo := &BidEntityMongo{
				ID:        bidValue.ID,
				AuctionID: bidValue.AuctionID,
				UserID:    bidValue.UserID,
				Amount:    bidValue.Amount,
				Timestamp: bidValue.Timestamp.Unix(),
			}

			bd.Collection.InsertOne(ctx, bidEntityMongo)

		}(bid)

	}

	wg.Wait()

	return nil
}
