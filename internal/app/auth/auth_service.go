package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	AuthDTO "github.com/edwinedjokpa/event-booking-api/internal/app/auth/dto"
	"github.com/edwinedjokpa/event-booking-api/internal/app/user"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/service/otp"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/service/session"
	HTTPException "github.com/edwinedjokpa/event-booking-api/internal/pkg/shared/httpexception"
	"github.com/edwinedjokpa/event-booking-api/internal/pkg/util"

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
	sessionService *session.SessionService
	otpService     otp.OTPService
}

func NewAuthService(repository user.UserRepository, jwtSecret string, sessionService *session.SessionService, otpService otp.OTPService) AuthService {
	return &authService{repository, jwtSecret, sessionService, otpService}
}

func (svc *authService) Register(request AuthDTO.RegisterUserRequest) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	existingUser, dbErr := svc.repository.FindOneByEmail(normalizedEmail)

	if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		panic(dbErr)
	}

	if existingUser != nil {
		_, _ = util.HashPassword("dummy_password_for_security")
		panic(HTTPException.NewConflictException("User with email already exists", nil))
	}

	hashedPassword, err := util.HashPassword(request.Password)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to hash user password", err.Error()))
	}

	newUser := user.User{
		ID:        util.GenerateUUID(),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     normalizedEmail,
		Password:  string(hashedPassword),
	}

	if err := svc.repository.Create(newUser); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to create user account", err.Error()))
	}
}

func (svc *authService) Login(ctx context.Context, request AuthDTO.LoginUserRequest) AuthDTO.LoginResponse {
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	user, dbErr := svc.repository.FindOneByEmail(normalizedEmail)

	var storedPassword string
	if dbErr != nil {
		storedPassword = "dummy_hash_for_security"
	} else {
		storedPassword = user.Password
	}

	isValid := util.CheckPasswordHash(storedPassword, request.Password)

	isUserNotFound := errors.Is(dbErr, gorm.ErrRecordNotFound)
	if isUserNotFound || !isValid {
		panic(HTTPException.NewBadRequestException("Invalid credentials", nil))
	}

	if dbErr != nil {
		panic(dbErr)
	}

	accessExpiresAt := 1 * time.Hour
	accessClaims := jwt.MapClaims{"userID": user.ID, "email": user.Email}
	accessToken, err := util.GenerateToken(accessClaims, accessExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate access token", err.Error()))
	}

	refreshSessionID := util.GenerateUUID()
	refreshExpiresAt := 7 * 24 * time.Hour
	refreshClaims := jwt.MapClaims{"sessionID": refreshSessionID}
	refreshToken, err := util.GenerateToken(refreshClaims, refreshExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate refresh token", err.Error()))
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
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

	user, dbErr := svc.repository.FindOneByEmail(normalizedEmail)
	otpErr := svc.otpService.ValidateOTP(normalizedEmail, request.OTP)

	isUserNotFound := errors.Is(dbErr, gorm.ErrRecordNotFound)
	if isUserNotFound || otpErr != nil {
		panic(HTTPException.NewBadRequestException("Invalid or expired OTP", nil))
	}

	if dbErr != nil {
		panic(dbErr)
	}

	newHashedPassword, err := util.HashPassword(request.NewPassword)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to hash new password", err))
	}

	err = svc.repository.UpdatePassword(user.ID, newHashedPassword)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to update password", err))
	}
}

func (svc *authService) Logout(ctx context.Context, refreshToken string) {
	_, claims, err := util.ValidateToken(refreshToken, []byte(svc.jwtSecret))
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
	_, claims, err := util.ValidateToken(tokenString, []byte(svc.jwtSecret))
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

	newSessionID := util.GenerateUUID()
	refreshExpiresAt := 7 * 24 * time.Hour

	if err := svc.sessionService.SetSession(ctx, newSessionID, sessionData.UserID, sessionData.Email, refreshExpiresAt); err != nil {
		panic(HTTPException.NewBadRequestException("Failed to save new refresh session", nil))
	}

	accessExpiresAt := 1 * time.Hour
	accessClaims := jwt.MapClaims{"userID": sessionData.UserID, "email": sessionData.Email}
	newAccessToken, err := util.GenerateToken(accessClaims, accessExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate new access token", nil))
	}

	refreshClaims := jwt.MapClaims{"sessionID": newSessionID}
	newRefreshToken, err := util.GenerateToken(refreshClaims, refreshExpiresAt, svc.jwtSecret)
	if err != nil {
		panic(HTTPException.NewBadRequestException("Failed to generate new refresh token", nil))
	}

	return AuthDTO.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
}
