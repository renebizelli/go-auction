package auction_entity

import (
	"context"
	"renebizelli/go-leilao/Internals/internal_error"
	"time"

	"github.com/google/uuid"
)

type Auction struct {
	ID          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active    AuctionStatus = iota
	Completed               = 1
)

const (
	New         ProductCondition = iota
	Used                         = 1
	Refurbished                  = 2
)

func CreateAuction(productName, category, description string,
	condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := Auction{
		ID:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if v := auction.Validate(); v != nil {
		return nil, v
	}

	return &auction, nil
}

func (a *Auction) Validate() *internal_error.InternalError {
	if len(a.ProductName) <= 1 ||
		len(a.Category) < 2 ||
		len(a.Description) < 10 &&
			(a.Condition != New || a.Condition != Used || a.Condition != Refurbished) {
		return internal_error.NewBadRequestError("Invalid auction data")
	}

	return nil
}

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auctionEntity *Auction) *internal_error.InternalError
	FindAuctionByID(ctx context.Context, id string) (*Auction, *internal_error.InternalError)
	FindAuction(ctx context.Context, status AuctionStatus, category, productName string) ([]Auction, *internal_error.InternalError)
}
