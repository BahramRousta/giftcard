package adaptor

import (
	"errors"
	"fmt"
	"giftCard/util"
	"github.com/rs/zerolog/log"
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
		BaseUrl:      config.BASE_URL,
		ClientID:     config.CLIENT_ID,
		ClientSecret: config.CLIENT_SECRET,
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
