package giftcard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"giftCard/config"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type GiftCard struct {
	BaseUrl      string `json:"baseUrl"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	Logger       *zap.Logger
}

func NewGiftCard(logger *zap.Logger) *GiftCard {
	return &GiftCard{
		BaseUrl:      config.C().Service.BaseUrl,
		ClientID:     config.C().Service.ClientID,
		ClientSecret: config.C().Service.ClientSecret,
		Logger:       logger,
	}
}

func (g *GiftCard) ProcessRequest(ctx context.Context, method string, url string, payload *[]byte) (map[string]any, error) {

	uniqueID, _ := ctx.Value("tracer").(string)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	token, err := g.Auth(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token.Token)

	if payload != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(*payload))
	}

	logger := g.Logger.With(
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
		return nil, errors.New("error while sending request")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error while process response")
	}

	logger.Info("response from provider",
		zap.Int("status_code", res.StatusCode),
		zap.Any("headers", res.Header),
		zap.Any("body", res.Body),
	)

	if res.StatusCode == http.StatusForbidden {
		return nil, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, errors.New("error while unmarshal response body")
	}

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	}
	return responseData, &RequestErr{ErrMsg: "error from provider", Response: responseData}
}
