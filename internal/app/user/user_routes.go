package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller UserController) {
	router.POST("/dashboard", controller.Dashboard)
}
