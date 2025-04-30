package user_controller

import (
	"net/http"
	"renebizelli/go-leilao/Internals/infra/api/web/validations"
	"renebizelli/go-leilao/Internals/usecase/user_usecase"
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase user_usecase.UserUseCaseInterface
}

func NewUserController(userUseCase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (u *UserController) CreateUser(c *gin.Context) {

	var userInputDTO user_usecase.UserInputDTO

	if err := c.ShouldBindJSON(&userInputDTO); err != nil {
		restErr := validations.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	if err := u.userUseCase.CreateUser(c, userInputDTO); err != nil {
		restErr := rest_err.ConvertError(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
