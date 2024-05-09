package giftcard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"giftcard/config"
	"giftcard/internal/adaptor/redis"
	"giftcard/internal/adaptor/trace"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
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

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error while creating new request", err.Error()))
		return nil, err
	}

	token, err := g.Auth(spannedContext)
	if err != nil {
		span.SetAttributes(attribute.String("error while authentication to gift card provide", err.Error()))
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
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	var res *http.Response
	var bodyBytes []byte
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = client.Do(req)

		if err != nil {
			var urlErr *_url.Error
			if errors.As(err, &urlErr) {
				span.SetAttributes(attribute.String("URL error encountered", urlErr.Error()))
				time.Sleep(1 * time.Second)
				continue
			}

			if isEOFError(err) {
				span.SetAttributes(attribute.String("EOF error encountered", urlErr.Error()))
				time.Sleep(1 * time.Second)
				continue
			}

			span.SetAttributes(attribute.String("error while sending request to gift card provider", err.Error()))
			return nil, err
		}
		defer res.Body.Close()

		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			span.SetAttributes(attribute.String("error while reading response body", err.Error()))
			return nil, err
		}
		break
	}

	var responseData map[string]any
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, errors.New("error while unmarshal response body")
	}

	response := responser.GiftCardResponse{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Body:       string(bodyBytes),
	}
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	switch res.StatusCode {
	case http.StatusOK:
		return responseData, nil
	case http.StatusForbidden:
		return responseData, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	default:
		return responseData, &RequestErr{ErrMsg: "error from provider", Response: responseData}
	}
}

func isEOFError(err error) bool {
	return err != nil && err == io.EOF
}
