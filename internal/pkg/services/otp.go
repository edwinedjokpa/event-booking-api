package services

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thanhpk/randstr"
)

type OTPService interface {
	GenerateAndStoreOTP(email string) (string, error)
	ValidateOTP(email, otp string) error
}

type otpService struct {
	redisClient *redis.Client
}

func NewOTPService(redisClient *redis.Client) OTPService {
	return &otpService{
		redisClient: redisClient,
	}
}

func (s *otpService) set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return s.redisClient.Set(ctx, key, value, expiration).Err()
}

func (s *otpService) get(key string) (string, error) {
	ctx := context.Background()
	val, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *otpService) del(key string) error {
	ctx := context.Background()
	return s.redisClient.Del(ctx, key).Err()
}

func (s *otpService) GenerateAndStoreOTP(email string) (string, error) {
	otp := randstr.String(6, "0123456789")

	otpExpiresAt := 15 * time.Minute
	if err := s.set("otp_code:"+email, otp, otpExpiresAt); err != nil {
		return "", err
	}

	return otp, nil
}

func (s *otpService) ValidateOTP(email, userOTP string) error {
	key := "otp_code:" + email
	storedOTP, err := s.get(key)
	if err != nil {
		return err
	}

	if storedOTP != userOTP {
		return errors.New("invalid OTP")
	}

	s.del(key)
	return nil
}
