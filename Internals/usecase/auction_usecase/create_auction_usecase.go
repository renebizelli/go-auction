package auction_usecase

import (
	"context"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"renebizelli/go-leilao/Internals/entity/bid_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/Internals/usecase/bid_usecase"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition" binding:"oneof=1 2 3"`
}

type AuctionOutputDTO struct {
	ID          string           `json:"_id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02T15:04:05Z"`
}

type ProductCondition int64
type AuctionStatus int64

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

type AuctionUseCase struct {
	auctionRepository auction_entity.AuctionRepositoryInterface
	bidRepository     bid_entity.BidRepositoryInterface
}

type AuctionUseCaseInterface interface {
	CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) *internal_error.InternalError
	FindAuctionByID(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)
	FindWinningByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

func NewAuctionUseCase(auctionRepository auction_entity.AuctionRepositoryInterface, bidRepository bid_entity.BidRepositoryInterface) AuctionUseCaseInterface {
	return &AuctionUseCase{
		auctionRepository: auctionRepository,
		bidRepository:     bidRepository,
	}
}

func (a *AuctionUseCase) CreateAuction(ctx context.Context, auctionInput AuctionInputDTO) *internal_error.InternalError {

	auctionEntity, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition),
	)

	if err != nil {
		return err
	}

	if e := a.auctionRepository.CreateAuction(ctx, auctionEntity); e != nil {
		return e
	}

	return nil
}
