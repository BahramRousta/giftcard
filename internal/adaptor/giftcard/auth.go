package adaptor

import (
	"errors"
	rds "giftCard/internal/adaptor/redis"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

func (g *GiftCard) Auth() (string, error) {

	conn := rds.GetRedisConn()
	defer conn.Close()
	token, err := redis.String(conn.Do("GET", "giftcard_token"))

	if err == nil {
		return token, nil
	}

	url := g.BaseUrl + "/auth/jwt"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("client-id", g.ClientID)
	req.Header.Add("client-secret", g.ClientSecret)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		authHeader := res.Header.Get("Authorization")
		conn.Do("SET", "giftcard_token", authHeader, "EX", 3600)
		return authHeader, nil
	}
	return "", errors.New("authentication failed")
}
