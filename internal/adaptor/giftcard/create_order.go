package giftcard

import (
	"context"
	"encoding/json"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type ModifiedDate struct {
	Seconds     int64 `json:"_seconds"`
	Nanoseconds int   `json:"_nanoseconds"`
}

type exchangeRate struct {
	BaseCurrency   string       `json:"baseCurrency"`
	ModifiedDate   ModifiedDate `json:"modifiedDate"`
	Rate           float64      `json:"rate"`
	TargetCurrency string       `json:"targetCurrency"`
}

type Effect struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type MetaData struct {
	BaseCurrency   string       `json:"baseCurrency,omitempty"`
	ModifiedDate   ModifiedDate `json:"modifiedDate,omitempty"`
	Rate           float64      `json:"rate,omitempty"`
	TargetCurrency string       `json:"targetCurrency,omitempty"`
	Quantity       int          `json:"quantity,omitempty"`
	Quote          int          `json:"quote,omitempty"`
}

type Item struct {
	Description string   `json:"description"`
	Effect      Effect   `json:"effect"`
	MetaData    MetaData `json:"metaData"`
	Type        string   `json:"type"`
}

type Record struct {
	Items []Item `json:"items"`
	Sku   string `json:"sku"`
	Total struct {
		DKK float64 `json:"DKK"`
		EUR float64 `json:"EUR"`
	} `json:"total"`
}

type Invoice struct {
	PaymentMethod string   `json:"paymentMethod"`
	Records       []Record `json:"records"`
	Status        string   `json:"status"`
	Total         float64  `json:"total"`
	Wallet        string   `json:"wallet"`
}

type OrderResponse struct {
	Data struct {
		ExchangeRates []exchangeRate `json:"exchangeRates"`
		ExpiresAt     int64          `json:"expiresAt"`
		ID            string         `json:"id"`
		Invoice       Invoice        `json:"invoice"`
	} `json:"data"`
}

func (g *GiftCard) CreateOrder(ctx context.Context, productList []map[string]any) (OrderResponse, error) {

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

	url := g.BaseUrl + "/order/create"
	method := "POST"

	payload := map[string]any{
		"productList": productList,
		"wallet":      "EUR",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("error while prepare request to gift card provider",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return OrderResponse{}, err
	}

	data, err := g.ProcessRequest(spannedContext, method, url, &payloadBytes)
	if err != nil {
		logger.Error("error while processing request to gift card provider",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return OrderResponse{}, err
	}

	jsonData, err := json.Marshal(data)

	var responseData OrderResponse
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		logger.Error("error while unmarshal gift card provider response data",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return OrderResponse{}, err
	}
	return responseData, nil

}

type CreateOrderError struct {
	ErrorMsg string
	Response []byte
}

func (e *CreateOrderError) Error() string {
	return e.ErrorMsg
}
