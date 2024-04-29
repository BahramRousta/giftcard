package adaptor

import (
	"errors"
	"net/http"
)

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

	if res.StatusCode == http.StatusOK {
		authHeader := res.Header.Get("Authorization")
		return authHeader, nil
	}
	return "", errors.New("authentication failed")
}
