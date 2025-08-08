package middleware

import (
	"strings"

	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/util"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(
				HTTPException.NewUnauthorizedException("Authorization Header is missing", nil).StatusCode,
				HTTPException.NewUnauthorizedException("Authorization Header is missing", nil).ToResponse(),
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(
				HTTPException.NewUnauthorizedException("Invalid Authorization Header format", nil).StatusCode,
				HTTPException.NewUnauthorizedException("Invalid Authorization Header format", nil).ToResponse(),
			)
			return
		}

		tokenString := parts[1]

		token, claims, err := util.ValidateToken(tokenString, jwtKey)
		if err != nil || token == nil || !token.Valid {
			ctx.AbortWithStatusJSON(
				HTTPException.NewUnauthorizedException("Invalid token", nil).StatusCode,
				HTTPException.NewUnauthorizedException("Invalid token", nil).ToResponse(),
			)
			return
		}

		userID, userIDExists := claims["userID"].(string)
		if !userIDExists {
			ctx.AbortWithStatusJSON(
				HTTPException.NewUnauthorizedException("Token claims are missing required data", nil).StatusCode,
				HTTPException.NewUnauthorizedException("Token claims are missing required data", nil).ToResponse(),
			)
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
