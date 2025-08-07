package services

import (
	"context"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
)

type OTPService interface {
	GenerateAndStoreOTP(email string) (string, error)
	ValidateOTP(email, otp string) error
	DeleteOTPKey(email string) error
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
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "My App",
		AccountName: email,
	})

	if err != nil {
		return "", err
	}

	otpExpiresAt := 15 * time.Minute
	if err := s.set("otp_key:"+email, key.Secret(), otpExpiresAt); err != nil {
		return "", err
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (s *otpService) ValidateOTP(email, otp string) error {
	otpKey, err := s.get("otp_key:" + email)
	if err != nil {
		return err
	}

	if !totp.Validate(otp, otpKey) {
		return err
	}

	return nil
}

func (s *otpService) DeleteOTPKey(email string) error {
	return s.del("otp_key:" + email)
}
