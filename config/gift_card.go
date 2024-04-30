package config

type GiftCard struct {
	BaseUrl      string `mapstructure:"base_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}
