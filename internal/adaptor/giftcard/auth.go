package giftcard

import (
	"context"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
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
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	var authToken AuthToken
	getTokenErr := g.redis.Get(spannedContext, "giftcard_token", &authToken.Token)
	if getTokenErr != nil {
		logger.Error("internal error", zap.String("message", getTokenErr.Error()))
		span.SetAttributes(attribute.String("internal error", getTokenErr.Error()))
	}

	if authToken.Token != "" {
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
	logger.Error("Request to provider", zap.Any("message", request))
	span.SetAttributes(attribute.String("Request to provider", utils.Marshal(request)))

	if err != nil {
		return AuthToken{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("internal error", zap.String("error", err.Error()))
		span.SetAttributes(attribute.String("internal error", err.Error()))
		return AuthToken{}, &InternalErr{ErrMsg: err.Error()}
	}

	if res.StatusCode == http.StatusForbidden {
		logger.Error("Response from provider", zap.String("data", exceptions.StatusForbidden))
		span.SetAttributes(attribute.String("Response from provider", exceptions.StatusForbidden))
		return AuthToken{}, &ForbiddenErr{ErrMsg: exceptions.StatusForbidden}
	}

	response := responser.GiftCardResponse{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       string(bodyBytes),
	}
	logger.Info("Response from provider", zap.Any("data", response))
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	switch res.StatusCode {
	case http.StatusOK:
		authHeader := res.Header.Get("Authorization")
		setToRedisErr := g.redis.Set(spannedContext, "giftcard_token", authHeader, 3600*time.Second)
		if setToRedisErr != nil {
			logger.Error("internal error", zap.String("error", setToRedisErr.Error()))
			span.SetAttributes(attribute.String("internal error", setToRedisErr.Error()))
		}
		authToken.Token = authHeader
		return authToken, nil
	default:
		span.SetAttributes(attribute.String("error", exceptions.AuthenticationError),
			attribute.String("error", string(bodyBytes)),
		)
		return AuthToken{}, &AuthErr{ErrMsg: exceptions.AuthenticationError}
	}

}
