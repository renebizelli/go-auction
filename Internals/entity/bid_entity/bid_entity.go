package bid_entity

import (
	"context"
	"renebizelli/go-leilao/Internals/internal_error"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        string
	AuctionID string
	UserID    string
	Amount    float64
	Timestamp time.Time
}

func CreateBid(auctionId, userId string, amount float64) (*Bid, *internal_error.InternalError) {

	bid := &Bid{
		ID:        uuid.NewString(),
		AuctionID: auctionId,
		UserID:    userId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validate(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) Validate() *internal_error.InternalError {

	if e := uuid.Validate(b.AuctionID); e != nil {
		return internal_error.NewBadRequestError("Invalid AuctionID")
	}

	if e := uuid.Validate(b.UserID); e != nil {
		return internal_error.NewBadRequestError("Invalid UserID")
	}

	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Invalid Amount")
	}

	return nil
}

type BidRepositoryInterface interface {
	CreateBid(ctx context.Context, bids []Bid) *internal_error.InternalError
	FindByAuctionId(ctx context.Context, auctionId string) ([]Bid, *internal_error.InternalError)
	FindWinning(ctx context.Context, auctionId string) (*Bid, *internal_error.InternalError)
}
