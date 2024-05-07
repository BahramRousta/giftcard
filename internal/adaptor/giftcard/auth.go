package giftcard

import (
	"context"
	"errors"
	rds "giftcard/internal/adaptor/redis"
	"giftcard/internal/adaptor/trace"
	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type AuthToken struct {
	Token string
}

func (g *GiftCard) Auth(ctx context.Context) (AuthToken, error) {
	span, _ := trace.T.SpanFromContext(
		ctx,
		"AuthGiftCardRequest",
		"adapter")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

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

	logger.Info("attempting authentication request to gift card provider",
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

	logger.Info("authentication response from gift card provider",
		zap.Int("status_code", res.StatusCode),
		zap.Any("headers", res.Header),
		zap.Any("body", string(bodyBytes)),
	)

	if res.StatusCode == http.StatusForbidden {
		logger.Error("forbidden to access auth end point", zap.Int("status_code", http.StatusForbidden))
		span.SetAttributes(attribute.String("error", "Forbidden to access auth end point."))
		return AuthToken{}, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	if res.StatusCode == http.StatusOK {
		authHeader := res.Header.Get("Authorization")
		span.SetAttributes(attribute.String("token", authHeader))
		_, err = conn.Do("SET", "giftcard_token", authHeader, "EX", 3600)
		if err != nil {
			logger.Error("can not set token into cache", zap.String("err", err.Error()))
		}
		authToken.Token = authHeader
		return authToken, nil
	}

	span.SetAttributes(attribute.String("error", "error while attempting authentication from provider"),
		attribute.String("error", string(bodyBytes)),
	)
	return AuthToken{}, &AuthErr{ErrMsg: "Authentication Failed"}
}
