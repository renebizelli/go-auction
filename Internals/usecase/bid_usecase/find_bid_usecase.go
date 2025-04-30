package bid_usecase

import (
	"context"
	"renebizelli/go-leilao/Internals/internal_error"
)

func (bd *BidUseCase) FindByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError) {

	bids, err := bd.BidRepository.FindByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	var bidOutputDTOs []BidOutputDTO

	for _, b := range bids {
		bidOutputDTOs = append(bidOutputDTOs, BidOutputDTO{
			ID:        b.ID,
			AuctionID: b.AuctionID,
			UserID:    b.UserID,
			Amount:    b.Amount,
			Timestamp: b.Timestamp.Unix(),
		})
	}

	return bidOutputDTOs, nil
}

func (bd *BidUseCase) FindWinningByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError) {

	bid, err := bd.BidRepository.FindWinning(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	dto := &BidOutputDTO{
		ID:        bid.ID,
		AuctionID: bid.AuctionID,
		UserID:    bid.UserID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp.Unix(),
	}

	return dto, nil
}
