package adaptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"giftCard/util"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type GiftCard struct {
	BaseUrl      string
	ClientID     string
	ClientSecret string
}

func NewGiftCard() *GiftCard {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	return &GiftCard{
		BaseUrl:      config.BaseUrl,
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
	}
}

func (g *GiftCard) Auth() (string, error) {
	url := "https://sandbox-api.core.hub.gift/auth/jwt"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("client-id", g.ClientID)
	req.Header.Add("client-secret", g.ClientSecret)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	if res.StatusCode == http.StatusOK {
		authHeader := res.Header.Get("Authorization")
		return authHeader, nil
	}
	return "", errors.New("authentication failed")
}

func (g *GiftCard) CustomerInfo() (map[string]any, error) {
	url := "https://sandbox-api.core.hub.gift/customer/info"
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
	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var responseData map[string]any

		err = json.Unmarshal(bodyBytes, &responseData)
		if err != nil {
			return nil, err
		}

		return responseData, nil
	}
	return nil, errors.New("failed to get customer info")
}
