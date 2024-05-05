package giftcard

import (
	"context"
	"errors"
	rds "giftCard/internal/adaptor/redis"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type AuthToken struct {
	Token string
}

func (g *GiftCard) Auth(ctx context.Context) (AuthToken, error) {

	uniqueID, _ := ctx.Value("tracer").(string)

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

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)
	logger.Info("attempting authentication request",
		zap.String("url", g.BaseUrl+"/auth/jwt"),
		zap.String("method", "GET"),
		zap.Any("body", req.Body),
		zap.Any("headers", req.Header),
		zap.Any("params", req.URL.Query()),
	)

	res, err := client.Do(req)
	if err != nil {
		return AuthToken{}, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return AuthToken{}, errors.New("error while process response")
	}

	logger.Info("response from provider",
		zap.Int("status_code", res.StatusCode),
		zap.Any("headers", res.Header),
		zap.Any("body", string(bodyBytes)),
	)

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
