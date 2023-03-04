package token_process

import (
	"github.com/go-redis/redis/v9"
)

func RemoveToken(token string) error {
	// connect redis
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	err := client.Del(ctx, token).Err()
		if err != nil {
			return err
		}

	return err
}