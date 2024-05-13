package giftcard

import (
	"context"
	"encoding/json"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type ExchangeRate struct {
	BaseCurrency string `json:"baseCurrency"`
	ModifiedDate struct {
		Nanoseconds int64 `json:"_nanoseconds"`
		Seconds     int64 `json:"_seconds"`
	} `json:"modifiedDate"`
	Rate           float64 `json:"rate"`
	TargetCurrency string  `json:"targetCurrency"`
}

type Wallet struct {
	EUR struct {
		Balance       float64 `json:"balance"`
		CreditBalance int     `json:"creditBalance"`
		FrozenBalance int     `json:"frozenBalance"`
	} `json:"EUR"`
}

type CustomerInfoResponse struct {
	Data struct {
		ExchangeRates []ExchangeRate `json:"exchangeRates"`
		Name          string         `json:"name"`
		Wallet        Wallet         `json:"wallet"`
	} `json:"data"`
}

func (g *GiftCard) CustomerInfo(ctx context.Context) (CustomerInfoResponse, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"CustomerInfoAdapter",
		"GiftCardAdapter",
	)
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	url := g.BaseUrl + "/customer/info"
	method := "GET"

	data, err := g.ProcessRequest(spannedContext, method, url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return CustomerInfoResponse{}, err
	}

	jsonData, err := json.Marshal(data)

	var responseData CustomerInfoResponse
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		logger.Error(exceptions.InternalServerError,
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return CustomerInfoResponse{}, err
	}

	span.SetAttributes(attribute.String("data", string(jsonData)))
	return responseData, nil
}

type CustomerInfoError struct {
	ErrorMsg   string
	StatusCode int
	Response   []byte
}

func (e *CustomerInfoError) Error() string {
	return e.ErrorMsg
}
