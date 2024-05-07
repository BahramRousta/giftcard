package giftcard

import (
	"context"
	"encoding/json"
	"fmt"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type CustomerDeal struct {
	Discount struct {
		AdjustmentMode string `json:"adjustmentMode"`
		Amount         int    `json:"amount"`
	} `json:"discount"`
	Fee struct {
		AdjustmentMode string `json:"adjustmentMode"`
		Amount         int    `json:"amount"`
	} `json:"fee"`
}

type Variant struct {
	Name  string `json:"name"`
	Quote int    `json:"quote"`
	SKU   string `json:"sku"`
	State string `json:"state"`
}

type ProductResponse struct {
	Data struct {
		BaseCurrency string             `json:"baseCurrency"`
		Country      string             `json:"country"`
		CustomerDeal CustomerDeal       `json:"customerDeal"`
		Name         string             `json:"name"`
		ParentID     string             `json:"parentId"`
		ProductID    string             `json:"productId"`
		ProductType  string             `json:"productType"`
		Region       string             `json:"region"`
		Variants     map[string]Variant `json:"variants"`
	} `json:"data"`
}

func (g *GiftCard) ShopItem(ctx context.Context, productId string) (ProductResponse, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"ShopItemAdapter",
		"GiftCardAdapter",
	)
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	url := g.BaseUrl + fmt.Sprintf("/shop/products/%s", productId)
	method := "GET"

	data, err := g.ProcessRequest(spannedContext, method, url, nil)
	if err != nil {
		logger.Error("error while processing request to gift card provider",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return ProductResponse{}, err
	}

	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))

	var responseData ProductResponse
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return ProductResponse{}, err
	}
	return responseData, nil
}
