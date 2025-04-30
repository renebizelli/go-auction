package bid_usecase

import (
	"context"
	"os"
	"renebizelli/go-leilao/Internals/entity/bid_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/configuration/logger"
	"strconv"
	"time"
)

type BidInputDTO struct {
	AuctionID string  `json:"auction_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	ID        string  `json:"_id"`
	AuctionID string  `json:"auction_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp" time_format:"2006-01-02T15:04:05Z07:00"`
}

type BidUseCase struct {
	BidRepository       bid_entity.BidRepositoryInterface
	time                time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidCH               chan bid_entity.Bid
}

func NewBidUseCase(bidRepository bid_entity.BidRepositoryInterface) BidUseCaseInterface {

	batchInsertInterval := getBatchInsertInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		time:                *time.NewTimer(batchInsertInterval),
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: batchInsertInterval,
		bidCH:               make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateRouting(context.Background())

	return bidUseCase
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bid BidInputDTO) *internal_error.InternalError
	FindByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
}

var bidBatch []bid_entity.Bid

func (u *BidUseCase) triggerCreateRouting(ctx context.Context) {

	go func() {
		defer close(u.bidCH)

		for {

			u.time = *time.NewTimer(u.batchInsertInterval)

			select {
			case bidentity, ok := <-u.bidCH:
				if !ok && len(bidBatch) > 0 {
					if e := u.BidRepository.CreateBid(ctx, bidBatch); e != nil {
						logger.Error("Error creating bid batch: ", e)
					}
					return
				}

				bidBatch = append(bidBatch, bidentity)

				if len(bidBatch) >= u.maxBatchSize {
					if e := u.BidRepository.CreateBid(ctx, bidBatch); e != nil {
						logger.Error("Error creating bid batch: ", e)
					}
					bidBatch = nil
				}

			case <-u.time.C:

				if e := u.BidRepository.CreateBid(ctx, bidBatch); e != nil {
					logger.Error("Error creating bid batch: ", e)
				}

				bidBatch = nil
			}
		}

	}()
}

func (u *BidUseCase) CreateBid(ctx context.Context, bid BidInputDTO) *internal_error.InternalError {

	entity, err := bid_entity.CreateBid(bid.AuctionID, bid.UserID, bid.Amount)

	if err != nil {
		return err
	}

	u.bidCH <- *entity

	return nil
}

func getBatchInsertInterval() time.Duration {

	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")

	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}
	return duration

}

func getMaxBatchSize() int {

	maxBatchSize, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))

	if err != nil {
		return 5
	}

	return maxBatchSize

}
