package adaptor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (g *GiftCard) ShopItem(productId string) (ProductResponse, error) {
	url := g.BaseUrl + fmt.Sprintf("/shop/products/%s", productId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ProductResponse{}, err
	}
	token, err := g.Auth()
	if err != nil {
		return ProductResponse{}, err
	}
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return ProductResponse{}, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ProductResponse{}, err
	}
	fmt.Println("stringbody", string(bodyBytes))
	if res.StatusCode != http.StatusOK {
		return ProductResponse{}, &ShopItemError{ErrorMsg: "failed to fetch shop item", Response: bodyBytes}
	}

	var responseData ProductResponse
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return ProductResponse{}, err
	}
	fmt.Println("responseData", responseData)
	return responseData, nil
}

type ShopItemError struct {
	ErrorMsg string
	Response []byte
}

func (e *ShopItemError) Error() string {
	return e.ErrorMsg
}
