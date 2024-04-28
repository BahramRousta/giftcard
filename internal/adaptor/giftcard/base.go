package adaptor

import (
	"giftCard/config"
	"github.com/rs/zerolog/log"
)

type GiftCard struct {
	BaseUrl      string
	ClientID     string
	ClientSecret string
}

func NewGiftCard() *GiftCard {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	return &GiftCard{
		BaseUrl:      config.BaseUrl,
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
	}
}
