package td_ameritrade_client_golang

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func GetClient(authCode string, consumerKey string, refreshCode string) (*http.Client, string, time.Time, error) {
	config := &oauth2.Config{
		ClientID:     consumerKey,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: urlAccessToken,
		},
	}

	restoredToken := &oauth2.Token{
		AccessToken:  authCode,
		RefreshToken: refreshCode,
		Expiry:       time.Now(),
		TokenType:    "Bearer",
	}
	tokenSource := config.TokenSource(oauth2.NoContext, restoredToken)
	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	savedToken, err := tokenSource.Token()
	handleFatalErr("Client/SavedToken", err)
	return client, savedToken.AccessToken, savedToken.Expiry, err
}
