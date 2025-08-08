package middleware

import (
	"log"
	"net/http"

	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)

				if httpErr, ok := r.(*HTTPException.HTTPException); ok {
					ctx.JSON(httpErr.StatusCode, httpErr.ToResponse())
					return
				}

				exception := HTTPException.NewInternalServerException(nil)
				ctx.JSON(http.StatusInternalServerError, exception.ToResponse())
			}
		}()
		ctx.Next()
	}
}
