package adaptor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g *GiftCard) ShopList(pageSize int, pageToken string) (map[string]any, error) {
	url := fmt.Sprintf("%s/shop/products?pageSize=%d&pageToken=%s", g.BaseUrl, pageSize, pageToken)
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
		return nil, &ShopListError{ErrorMsg: "failed to fetch shop list", Response: responseData}
	}
}

type ShopListError struct {
	ErrorMsg string
	Response map[string]any
}

func (e *ShopListError) Error() string {
	return e.ErrorMsg
}
