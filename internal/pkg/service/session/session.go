package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
}

type SessionService struct {
	client *redis.Client
}

func NewSessionService(client *redis.Client) *SessionService {
	return &SessionService{client: client}
}

func (s *SessionService) SetSession(ctx context.Context, sessionID, userID, email string, expiresAt time.Duration) error {
	sessionData := SessionData{
		UserID: userID,
		Email:  email,
	}

	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, sessionID, sessionJSON, expiresAt).Err()
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	sessionJSON, err := s.client.Get(ctx, sessionID).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	err = json.Unmarshal([]byte(sessionJSON), &sessionData)
	if err != nil {
		return nil, err
	}
	return &sessionData, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.client.Del(ctx, sessionID).Err()
}
