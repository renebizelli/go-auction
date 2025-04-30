package auction_usecase

import (
	"context"
	"fmt"
	"renebizelli/go-leilao/Internals/entity/auction_entity"
	"renebizelli/go-leilao/Internals/internal_error"
	"renebizelli/go-leilao/Internals/usecase/bid_usecase"
)

func (a *AuctionUseCase) FindAuctionByID(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {

	entity, err := a.auctionRepository.FindAuctionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := &AuctionOutputDTO{
		ID:          entity.ID,
		ProductName: entity.ProductName,
		Category:    entity.Category,
		Description: entity.Description,
		Condition:   ProductCondition(entity.Condition),
		Status:      AuctionStatus(entity.Status),
		Timestamp:   entity.Timestamp,
	}

	return dto, nil

}

func (a *AuctionUseCase) FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {

	entities, err := a.auctionRepository.FindAuction(ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	dtos := []AuctionOutputDTO{}

	for _, entity := range entities {

		dto := &AuctionOutputDTO{
			ID:          entity.ID,
			ProductName: entity.ProductName,
			Category:    entity.Category,
			Description: entity.Description,
			Condition:   ProductCondition(entity.Condition),
			Status:      AuctionStatus(entity.Status),
			Timestamp:   entity.Timestamp,
		}

		dtos = append(dtos, *dto)

	}

	return dtos, nil

}

func (bd *AuctionUseCase) FindWinningByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {

	fmt.Println("Auction ID 2:", auctionId)

	auction, err := bd.auctionRepository.FindAuctionByID(ctx, auctionId)
	if err != nil {
		fmt.Println("Erro 3", err.Error())

		return nil, err
	}

	fmt.Println("Auction ID 2 ProductName:", auction.ProductName)

	auctionDTO := &AuctionOutputDTO{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		Timestamp:   auction.Timestamp,
	}

	winningInfoOutputDTO := &WinningInfoOutputDTO{
		Auction: *auctionDTO,
		Bid:     nil,
	}

	fmt.Println("Auction ID 2 win:", auction.ProductName)

	bid, err := bd.bidRepository.FindWinning(ctx, auctionId)
	if err != nil {
		fmt.Println("erro: ", err.Error())

		return winningInfoOutputDTO, nil
	}

	fmt.Println("Amount: ", bid.Amount)

	winningInfoOutputDTO.Bid = &bid_usecase.BidOutputDTO{
		ID:        bid.ID,
		AuctionID: bid.AuctionID,
		UserID:    bid.UserID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp.Unix(),
	}

	return winningInfoOutputDTO, nil
}
