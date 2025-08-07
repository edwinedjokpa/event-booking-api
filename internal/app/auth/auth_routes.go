package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller AuthController) {
	router.POST("/auth/register", controller.Register)
	router.POST("/auth/login", controller.Login)
	router.POST("/auth/logout", controller.Logout)
	router.POST("/auth/refresh", controller.RefreshToken)
}
