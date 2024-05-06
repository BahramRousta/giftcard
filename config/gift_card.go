package config

type GiftCard struct {
	Name         string `mapstructure:"name"`
	BaseUrl      string `mapstructure:"base_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}
