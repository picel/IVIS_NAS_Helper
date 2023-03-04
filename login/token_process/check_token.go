package token_process

import (
	"context"
	"fmt"
	"ivis_nas/key"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func VerifyToken(r *http.Request) (string, error) {
	tokenString, err := r.Cookie("jwt_token")
	if err != nil || tokenString.Value == "" {
		err = fmt.Errorf("token not found")
		return "", err
	}

	// verify token
	token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return key.JWTVerifyKey, nil
	})

	// connect redis
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	// get uuid from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return "", err
		}

		// get userid from redis
		userid, err := client.Get(ctx, accessUuid).Result()
		if err != nil {
			return "", err
		}

		// set ttl 900s to redis
		err = client.Expire(ctx, accessUuid, 15*time.Minute).Err()

		return userid, nil
	}

	return "", err
}
