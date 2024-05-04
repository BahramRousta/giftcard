package giftcard

import (
	"bytes"
	"encoding/json"
	"errors"
	"giftCard/config"
	"io"
	"net/http"
)

type GiftCard struct {
	BaseUrl      string `json:"baseUrl"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func NewGiftCard() *GiftCard {
	return &GiftCard{
		BaseUrl:      config.C().Service.BaseUrl,
		ClientID:     config.C().Service.ClientID,
		ClientSecret: config.C().Service.ClientSecret,
	}
}

func (g *GiftCard) ProcessRequest(method string, url string, payload *[]byte) (map[string]any, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	token, err := g.Auth()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token.Token)

	if payload != nil {
		req.Body = io.NopCloser(bytes.NewBuffer(*payload))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error while sending request")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error while process response")
	}

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
