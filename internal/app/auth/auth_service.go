package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	AuthDTO "github.com/edwinedjokpa/event-booking-api/internal/app/auth/dto"
	"github.com/edwinedjokpa/event-booking-api/internal/app/user"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/services"
	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(request AuthDTO.RegisterUserRequest)
	Login(ctx context.Context, request AuthDTO.LoginUserRequest) AuthDTO.LoginResponse
	ForgotPassword(request AuthDTO.ForgotPasswordRequest)
	ResetPassword(request AuthDTO.ResetPasswordRequest)
	Logout(ctx context.Context, token string)
	RefreshToken(ctx context.Context, token string) AuthDTO.LoginResponse
}

type authService struct {
	repository     user.UserRepository
	jwtSecret      string
	sessionService *services.SessionService
	otpService     services.OTPService
}

func NewAuthService(repository user.UserRepository, jwtSecret string, sessionService *services.SessionService, otpService services.OTPService) AuthService {
	return &authService{repository, jwtSecret, sessionService, otpService}
}

func (svc *authService) Register(request AuthDTO.RegisterUserRequest) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	existingUser, err := svc.repository.FindOneByEmail(normalizedEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if existingUser != nil {
		panic(HTTPException.NewConflictException("User with email already exists", existingUser.Email))
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to hash user password", err.Error()))
	}

	user := user.User{
		ID:        utils.GenerateUUID(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     normalizedEmail,
		Password:  string(hashedPassword),
	}

	if err := svc.repository.Create(user); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to create user account", err.Error()))
	}
}

func (svc *authService) Login(ctx context.Context, request AuthDTO.LoginUserRequest) AuthDTO.LoginResponse {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := svc.repository.FindOneByEmail(normalizedEmail)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Invalid credentials", nil))
	}

	isValid := utils.CheckPasswordHash(user.Password, request.Password)
	if !isValid {
		panic(HTTPException.NewBadRequestException("Invalid credentials", nil))
	}

	accessExpiresAt := 1 * time.Hour
	accessClaims := jwt.MapClaims{"userID": user.ID, "email": user.Email}
	accessToken, err := utils.GenerateToken(accessClaims, accessExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate access token", nil))
	}

	refreshSessionID := utils.GenerateUUID()
	refreshExpiresAt := 7 * 24 * time.Hour
	refreshClaims := jwt.MapClaims{"sessionID": refreshSessionID}
	refreshToken, err := utils.GenerateToken(refreshClaims, refreshExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate refresh token", nil))
	}

	err = svc.sessionService.SetSession(ctx, refreshSessionID, user.ID, user.Email, refreshExpiresAt)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to create session", err.Error()))
	}

	return AuthDTO.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (svc *authService) ForgotPassword(request AuthDTO.ForgotPasswordRequest) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := svc.repository.FindOneByEmail(normalizedEmail)
	if err != nil {
		return
	}

	otp, err := svc.otpService.GenerateAndStoreOTP(user.Email)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate OTP", err.Error()))
	}

	fmt.Printf("OTP for %s is: %s (expires in 15 minutes)\n", user.Email, otp)
}

func (svc *authService) ResetPassword(request AuthDTO.ResetPasswordRequest) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	user, err := svc.repository.FindOneByEmail(normalizedEmail)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Invalid or expired OTP", nil))
	}

	err = svc.otpService.ValidateOTP(user.Email, request.OTP)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Invalid or expired OTP", nil))
	}

	newHashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to hash new password", nil))
	}

	err = svc.repository.UpdatePassword(user.ID, newHashedPassword)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to update password", err.Error()))
	}

	err = svc.otpService.DeleteOTPKey(user.Email)
	if err != nil {
		return
	}
}

func (svc *authService) Logout(ctx context.Context, refreshToken string) {
	_, claims, err := utils.ValidateToken(refreshToken, []byte(svc.jwtSecret))
	if err != nil {
		panic(HTTPException.NewUnauthorizedException("Invalid refresh token", nil))
	}

	sessionID, ok := claims["sessionID"].(string)
	if !ok {
		panic(HTTPException.NewUnauthorizedException("invalid session ID in refresh token claims", nil))
	}

	err = svc.sessionService.DeleteSession(ctx, sessionID)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to delete session", err.Error()))
	}
}

func (svc *authService) RefreshToken(ctx context.Context, tokenString string) AuthDTO.LoginResponse {
	_, claims, err := utils.ValidateToken(tokenString, []byte(svc.jwtSecret))
	if err != nil {
		panic(HTTPException.NewUnauthorizedException("Invalid refresh token", nil))
	}

	sessionIDRaw, exists := claims["sessionID"]
	if !exists {
		panic(HTTPException.NewBadRequestException("Invalid session", nil))
	}

	sessionID, ok := sessionIDRaw.(string)
	if !ok {
		panic(HTTPException.NewUnauthorizedException("Invalid session", nil))
	}

	sessionData, err := svc.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		panic(HTTPException.NewUnauthorizedException("Session expired or revoked", nil))
	}

	if err := svc.sessionService.DeleteSession(ctx, sessionID); err != nil {
		panic(HTTPException.NewBadRequestException("failed to delete sessionI D", nil))
	}

	newSessionID := utils.GenerateUUID()
	refreshExpiresAt := 7 * 24 * time.Hour

	if err := svc.sessionService.SetSession(ctx, newSessionID, sessionData.UserID, sessionData.Email, refreshExpiresAt); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to save new refresh session", nil))
	}

	accessExpiresAt := 15 * time.Minute
	accessClaims := jwt.MapClaims{"userID": sessionData.UserID, "email": sessionData.Email}
	newAccessToken, err := utils.GenerateToken(accessClaims, accessExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate new access token", nil))
	}

	refreshClaims := jwt.MapClaims{"sessionID": newSessionID}
	newRefreshToken, err := utils.GenerateToken(refreshClaims, refreshExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate new refresh token", nil))
	}

	return AuthDTO.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
}
