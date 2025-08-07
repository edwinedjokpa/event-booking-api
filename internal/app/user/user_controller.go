package user

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Dashboard(c *gin.Context)
}

type userController struct {
	service UserService
}

func NewUserController(service UserService) UserController {
	return &userController{service}
}

func (ctrl *userController) Dashboard(ctx *gin.Context) {}
