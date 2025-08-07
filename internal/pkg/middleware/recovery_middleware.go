package middleware

import (
	"log"
	"net/http"

	HttpException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)

				if httpErr, ok := r.(*HttpException.ApiException); ok {
					ctx.JSON(httpErr.StatusCode, httpErr.ToResponse())
					return
				}

				exception := HttpException.NewInternalServerException(nil)
				ctx.JSON(http.StatusInternalServerError, exception.ToResponse())
			}
		}()
		ctx.Next()
	}
}
