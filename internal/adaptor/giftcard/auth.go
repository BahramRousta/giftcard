package giftcard

import (
	"context"
	"errors"
	"giftcard/internal/adaptor/trace"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"io"
	"net/http"
	"time"
)

type AuthToken struct {
	Token string
}

func (g *GiftCard) Auth(ctx context.Context) (AuthToken, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"AuthGiftCardRequest",
		"adapter")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	var authToken AuthToken
	getTokenErr := g.redis.Get(spannedContext, "giftcard_token", &authToken.Token)
	if getTokenErr != nil {
		span.SetAttributes(attribute.String("error while getting token from redis", getTokenErr.Error()))
	}

	if authToken.Token != "" {
		span.SetAttributes(attribute.String("get token from redis", authToken.Token))
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

	request := requester.Request{
		ID:          uniqueID,
		RequestBody: req.Body,
		Uri:         url,
		Method:      method,
		Header:      req.Header,
		Params:      req.URL.Query(),
	}
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	if err != nil {
		return AuthToken{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.SetAttributes(attribute.String("error while processing auth response", err.Error()))
		return AuthToken{}, errors.New("error while processing auth response")
	}

	response := responser.GiftCardResponse{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       string(bodyBytes),
	}
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	if res.StatusCode == http.StatusForbidden {
		return AuthToken{}, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	if res.StatusCode == http.StatusOK {
		authHeader := res.Header.Get("Authorization")

		setToRedisErr := g.redis.Set(spannedContext, "giftcard_token", authHeader, 3600*time.Second)
		if setToRedisErr != nil {
			span.SetAttributes(attribute.String("error while setting token to redis", setToRedisErr.Error()))
		}

		authToken.Token = authHeader
		return authToken, nil
	}

	span.SetAttributes(attribute.String("error", "error while attempting authentication from provider"),
		attribute.String("error", string(bodyBytes)),
	)
	return AuthToken{}, &AuthErr{ErrMsg: "Authentication Failed"}
}
