package auction_controller

import (
	"net/http"
	"renebizelli/go-leilao/Internals/infra/api/web/validations"
	"renebizelli/go-leilao/Internals/usecase/auction_usecase"
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (ac *AuctionController) CreateAuction(c *gin.Context) {

	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		restErr := validations.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	if err := ac.auctionUseCase.CreateAuction(c, auctionInputDTO); err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
