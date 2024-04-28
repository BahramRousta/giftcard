package adaptor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g *GiftCard) ShopItem(productId string) (map[string]any, error) {
	url := g.BaseUrl + fmt.Sprintf("/shop/products/%s", productId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	token, err := g.Auth()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	} else {
		return nil, &ShopItemError{ErrorMsg: "failed to fetch shop item", Response: responseData}
	}
}

type ShopItemError struct {
	ErrorMsg string
	Response map[string]any
}

func (e *ShopItemError) Error() string {
	return e.ErrorMsg
}
