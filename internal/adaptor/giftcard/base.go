package giftcard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"giftcard/config"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type GiftCard struct {
	BaseUrl      string `json:"baseUrl"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	//Logger       *zap.Logger
}

func NewGiftCard() *GiftCard {
	return &GiftCard{
		BaseUrl:      config.C().Service.BaseUrl,
		ClientID:     config.C().Service.ClientID,
		ClientSecret: config.C().Service.ClientSecret,
		//Logger:       logger,
	}
}

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
		logger.Error("error while creating new request",
			zap.String("err", err.Error()))
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	token, err := g.Auth(spannedContext)
	if err != nil {
		logger.Error("error while authentication to gift card provider",
			zap.String("err", err.Error()))
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token.Token)

	if payload != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(*payload))
	}

	logger.Info("request to provider",
		zap.String("url", url),
		zap.String("method", method),
		zap.Any("body", req.Body),
		zap.Any("headers", req.Header),
		zap.Any("params", req.URL.Query()),
	)

	res, err := client.Do(req)
	if err != nil {
		logger.Error("error while sending request to gift card provider",
			zap.String("err", err.Error()))
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, errors.New("error while sending request")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, errors.New("error while process response")
	}

	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, errors.New("error while unmarshal response body")
	}

	logger.Info("response from gift card provider",
		zap.Int("status_code", res.StatusCode),
		zap.Any("headers", res.Header),
		zap.Any("body", string(bodyBytes)),
	)

	if res.StatusCode == http.StatusForbidden {
		logger.Error("error while sending request to gift card provider", zap.Int("status_code", http.StatusForbidden))
		span.SetAttributes(attribute.String("error", "Forbidden to access end point."))
		return responseData, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	if res.StatusCode == http.StatusOK {
		span.SetAttributes(attribute.String("data", string(bodyBytes)))
		return responseData, nil
	}

	span.SetAttributes(attribute.String("error", "error from provider"),
		attribute.String("error", string(bodyBytes)),
		attribute.Int("status_code", res.StatusCode),
	)

	return responseData, &RequestErr{ErrMsg: "error from provider", Response: responseData}
}
