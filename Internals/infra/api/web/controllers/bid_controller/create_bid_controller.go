package bid_controller

import (
	"net/http"
	"renebizelli/go-leilao/Internals/infra/api/web/validations"
	"renebizelli/go-leilao/Internals/usecase/bid_usecase"
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin"
)

type BidController struct {
	bidUseCase bid_usecase.BidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (bc *BidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validations.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	if err := bc.bidUseCase.CreateBid(c, bidInputDTO); err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
