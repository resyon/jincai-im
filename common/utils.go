package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func GetUUID() string {
	u, _ := uuid.NewUUID()
	return u.String()
}

const (
	salt      = "fJeUjYFj*3:{*U"
	UserIdKey = "user_id"
)

func AddSalt(raw string) string {
	sha := sha256.New()
	add := []byte(raw + salt)
	sha.Write(add)
	return hex.EncodeToString(sha.Sum(nil))
}

func GetUserIdFromContext(c *gin.Context) int {
	_userId, _ := c.Get(UserIdKey)
	userId, _ := _userId.(int)
	return userId
}

func SubUtilReady(client *redis.Client, channel string) (firstMsg *redis.Message, sub *redis.PubSub, err error) {
	sub = client.Subscribe(context.TODO(), channel)
	iface, err := sub.Receive(context.TODO())
	if err != nil {
		// handle error
		return nil, nil, err
	}

	// Should be *Subscription, but others are possible if other actions have been
	// taken on sub since it was created.
	switch v := iface.(type) {
	case *redis.Subscription:
		// subscribe succeeded
	case *redis.Message:
		// received first message
		firstMsg = v
	case *redis.Pong:
		// pong received
	default:
		// handle error
		return nil, nil, err
	}
	return
}
