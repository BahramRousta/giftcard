package giftcard

import (
	"context"
	"encoding/json"
	"fmt"
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
	url := g.BaseUrl + fmt.Sprintf("/shop/products/%s", productId)
	method := "GET"

	data, err := g.ProcessRequest(ctx, method, url, nil)
	jsonData, err := json.Marshal(data)

	var responseData ProductResponse
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return ProductResponse{}, err
	}
	return responseData, nil
}
