package auth

import (
	"net/http"

	AuthDTO "github.com/edwinedjokpa/event-booking-api/internal/app/auth/dto"
	APIResponse "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/apiresponse"
	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authController struct {
	service AuthService
}

func NewAuthController(service AuthService) AuthController {
	return &authController{service}
}

func (ctrl *authController) Register(ctx *gin.Context) {
	var request AuthDTO.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := HTTPException.NewBadRequestException("Bad Request Exception", err.Error())
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		exception := utils.FormatValidationErrors(err)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	ctrl.service.Register(request)
	ctx.JSON(http.StatusCreated, APIResponse.Success("User account created successfully", nil))
}

func (ctrl *authController) Login(ctx *gin.Context) {
	var request AuthDTO.LoginUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := HTTPException.NewBadRequestException("Bad Request Exception", err.Error())
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		exception := utils.FormatValidationErrors(err)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	tokens := ctrl.service.Login(request)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Domain:   "",
		MaxAge:   60 * 60 * 24 * 7,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(ctx.Writer, cookie)
	ctx.JSON(http.StatusOK, APIResponse.Success("User login successfully", gin.H{"token": tokens.AccessToken, "type": "Bearer"}))
}

func (ctrl *authController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusOK, APIResponse.Success("User logged out successfully", nil))
		return
	}

	ctrl.service.Logout(refreshToken)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(ctx.Writer, cookie)
	ctx.JSON(http.StatusOK, APIResponse.Success("User logged out successfully", nil))
}

func (ctrl *authController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		exception := HTTPException.NewUnauthorizedException("Refresh token is missing", nil)
		ctx.JSON(exception.StatusCode, exception.ToResponse())
		return
	}

	tokens := ctrl.service.RefreshToken(refreshToken)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Domain:   "",
		MaxAge:   60 * 60 * 24 * 7,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(ctx.Writer, cookie)
	ctx.JSON(http.StatusOK, APIResponse.Success("Tokens refreshed successfully", gin.H{"token": tokens.AccessToken}))
}
