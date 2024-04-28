package adaptor

import (
	"giftCard/util"
	"github.com/rs/zerolog/log"
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
