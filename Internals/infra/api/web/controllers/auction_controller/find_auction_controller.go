package auction_controller

import (
	"fmt"
	"renebizelli/go-leilao/Internals/usecase/auction_usecase"
	"renebizelli/go-leilao/configuration/rest_err"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *AuctionController) FindAuctionById(ctx *gin.Context) {

	auctionId := ctx.Param("auctionId")

	if uuid.Validate(auctionId) != nil {
		restErr := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid auction ID",
		})

		ctx.JSON(restErr.Code, restErr)
		return
	}

	auction, err := c.auctionUseCase.FindAuctionByID(ctx, auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		ctx.JSON(restErr.Code, restErr)
		return
	}

	ctx.JSON(200, auction)

}

func (c *AuctionController) FindAuctions(ctx *gin.Context) {

	status := ctx.Query("status")
	category := ctx.Query("category")
	productName := ctx.Query("productName")

	if status == "" {
		status = "0"
	}

	statusEnum, errConv := strconv.Atoi(status)

	if errConv != nil {
		restErr := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "status",
			Message: "Invalid status",
		})
		ctx.JSON(restErr.Code, restErr)
		return
	}

	auctions, err := c.auctionUseCase.FindAuctions(ctx, auction_usecase.AuctionStatus(statusEnum), category, productName)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		ctx.JSON(restErr.Code, restErr)
		return
	}

	ctx.JSON(200, auctions)
}

func (c *AuctionController) FindWinningBidByAuctionId(ctx *gin.Context) {

	auctionId := ctx.Param("auctionId")

	if uuid.Validate(auctionId) != nil {
		restErr := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid auction ID",
		})

		ctx.JSON(restErr.Code, restErr)
		return
	}

	fmt.Println("Auction ID:", auctionId)

	auction, err := c.auctionUseCase.FindWinningByAuctionId(ctx, auctionId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		ctx.JSON(restErr.Code, restErr)
		return
	}

	ctx.JSON(200, auction)
}
