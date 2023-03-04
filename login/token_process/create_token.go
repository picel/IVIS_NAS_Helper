package token_process

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"

	key "ivis_nas/key"
)

func CreateToken(userid string) (string, error) {
	// create random string
	accessUuid := uuid.New().String()

	// write access token to redis
	err := WriteDB(userid, accessUuid)
	if err != nil {
		return "", err
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_uuid": accessUuid,
		"user_id":     userid,
	})

	// sign token
	tokenString, err := token.SignedString(key.JWTVerifyKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WriteDB(userid string, accessUuid string) error {
	// connect redis
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	// set access token to redis with expiration time
	err := client.Set(ctx, accessUuid, userid, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
