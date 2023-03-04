package login

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	key "ivis_nas/key"
	token_process "ivis_nas/login/token_process"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

func RemoteLogin(w http.ResponseWriter, r *http.Request) {
	// get id and password from post request
	id := r.FormValue("id")
	pw := r.FormValue("pw")

	// send post request to 192.168.195.1/loginCheck and get response code
	resp, err := http.PostForm("http://192.168.195.1/loginCheck", url.Values{"id": {id}, "pw": {pw}})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// if response code is 200, login success
	if resp.StatusCode == 200 {
		// create jwt token
		ts, err := token_process.CreateToken(id)
		if err != nil {
			fmt.Println(err)
			// error
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			// success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(ts))
			return
		}
	} else {
		// error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func RemoteTokenVerify(w http.ResponseWriter, r *http.Request) {
	// json post request, get token
	tokenString := r.FormValue("token")
	fmt.Println(tokenString)

	// verify token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return key.JWTVerifyKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// connect redis
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	// get uuid from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get userid from redis
		userid, err := client.Get(ctx, accessUuid).Result()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// set ttl 900s to redis
		err = client.Expire(ctx, accessUuid, 15*time.Minute).Err()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userid))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	return
}
