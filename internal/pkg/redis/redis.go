package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client

type SessionData struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
}

func InitRedis(addr string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")
}

func SetSession(sessionID string, userID string, email string, expiresAt time.Duration) error {
	sessionData := map[string]string{
		"userID": userID,
		"email":  email,
	}
	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, sessionID, sessionJSON, expiresAt).Err()
}

func GetSession(sessionID string) (*SessionData, error) {
	sessionJSON, err := RedisClient.Get(ctx, sessionID).Result()
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

func DeleteSession(sessionID string) error {
	return RedisClient.Del(ctx, sessionID).Err()
}
