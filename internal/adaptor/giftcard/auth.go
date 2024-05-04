package giftcard

import (
	rds "giftCard/internal/adaptor/redis"
	"github.com/gomodule/redigo/redis"
	"net/http"
)

type AuthToken struct {
	Token string
}

func (g *GiftCard) Auth() (AuthToken, error) {

	conn := rds.GetRedisConn()
	defer conn.Close()

	token, err := redis.String(conn.Do("GET", "giftcard_token"))

	var authToken AuthToken
	if err == nil {
		authToken.Token = token
		return authToken, nil
	}

	url := g.BaseUrl + "/auth/jwt"
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return AuthToken{}, err
	}

	req.Header.Add("client-id", g.ClientID)
	req.Header.Add("client-secret", g.ClientSecret)

	res, err := client.Do(req)
	if err != nil {
		return AuthToken{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return AuthToken{}, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	if res.StatusCode == http.StatusOK {

		authHeader := res.Header.Get("Authorization")
		conn.Do("SET", "giftcard_token", authHeader, "EX", 3600)

		authToken.Token = authHeader
		return authToken, nil
	}
	return AuthToken{}, &AuthErr{ErrMsg: "Authentication Failed"}
}
