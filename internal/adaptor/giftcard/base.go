package giftcard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"giftcard/config"
	"giftcard/internal/adaptor/redis"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"io"
	"net/http"
	_url "net/url"
	"time"
)

type GiftCard struct {
	BaseUrl      string `json:"baseUrl"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	redis        *redis.Store
}

func NewGiftCard(redisStore *redis.Store) *GiftCard {
	return &GiftCard{
		BaseUrl:      config.C().Service.BaseUrl,
		ClientID:     config.C().Service.ClientID,
		ClientSecret: config.C().Service.ClientSecret,
		redis:        redisStore,
	}
}

const maxRetries = 3

func (g *GiftCard) ProcessRequest(ctx context.Context, method string, url string, payload *[]byte) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"ProcessGiftCardRequest",
		"adapter")

	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.Error(exceptions.InternalServerError, zap.String("error", err.Error()))
		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return nil, err
	}

	token, err := g.Auth(spannedContext)
	if err != nil {
		logger.Error(exceptions.AuthenticationError, zap.Any("error", err.Error()))
		span.SetAttributes(attribute.String(exceptions.AuthenticationError, err.Error()))
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token.Token)

	var requestBody string
	if payload != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(*payload))
		requestBody = string(*payload)
	}

	request := requester.Request{
		ID:          uniqueID,
		RequestBody: requestBody,
		Uri:         url,
		Method:      method,
		Header:      req.Header,
		Params:      req.URL.Query(),
	}

	logger.Info("Request to provider", zap.Any("data", request))
	span.SetAttributes(attribute.String("Request to provider", utils.Marshal(request)))

	var res *http.Response
	var bodyBytes []byte
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = client.Do(req)

		if err != nil {
			var urlErr *_url.Error
			if errors.As(err, &urlErr) {
				logger.Error(exceptions.InternalServerError, zap.String("error", urlErr.Error()))
				span.SetAttributes(attribute.String(exceptions.InternalServerError, urlErr.Error()))
				time.Sleep(1 * time.Second)
				continue
			}

			if isEOFError(err) {
				logger.Error(exceptions.InternalServerError, zap.String("error", err.Error()))
				span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
				time.Sleep(1 * time.Second)
				continue
			}
			logger.Error(exceptions.InternalServerError, zap.String("error", err.Error()))
			span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
			return nil, err
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			logger.Error(exceptions.InternalServerError, zap.String("error", err.Error()))
			span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
			return nil, err
		}
		break
	}

	if res.StatusCode == http.StatusForbidden {
		logger.Error("Response from provider", zap.String("data", exceptions.StatusForbidden))
		span.SetAttributes(attribute.String("Response from provider", exceptions.StatusForbidden))
		return nil, &ForbiddenErr{ErrMsg: exceptions.StatusForbidden}
	}

	var responseData map[string]any
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		logger.Error(exceptions.InternalServerError, zap.String("error", err.Error()))
		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return nil, &InternalErr{ErrMsg: exceptions.InternalServerError}
	}

	response := responser.GiftCardResponse{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       string(bodyBytes),
	}
	logger.Info("Response from provider", zap.Any("data", response))
	span.SetAttributes(attribute.String("Response from provider", utils.Marshal(response)))

	switch res.StatusCode {
	case http.StatusOK:
		return responseData, nil
	default:
		return responseData, &RequestErr{ErrMsg: "error from provider", Response: responseData}
	}
}

func isEOFError(err error) bool {
	return err != nil && err == io.EOF
}
