package bid_controller

import (
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (bc *BidController) FindBidByAuctionId(ctx *gin.Context) {

	auctionId := ctx.Param("auctionId")

	if uuid.Validate(auctionId) != nil {
		restErr := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid auction ID",
		})

		ctx.JSON(restErr.Code, restErr)
		return
	}

	bidOutputs, err := bc.bidUseCase.FindByAuctionId(ctx, auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		ctx.JSON(restErr.Code, restErr)
		return
	}

	ctx.JSON(200, bidOutputs)

}
