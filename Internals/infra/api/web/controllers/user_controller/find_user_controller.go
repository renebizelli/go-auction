package user_controller

import (
	"renebizelli/go-leilao/configuration/rest_err"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *UserController) FindUserByID(c *gin.Context) {

	userId := c.Param("userId")

	if uuid.Validate(userId) != nil {
		restErr := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "userId",
			Message: "Invalid user ID",
		})

		c.JSON(restErr.Code, restErr)
		return
	}

	user, err := u.userUseCase.FindUserByID(c, userId)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.JSON(200, user)
}
