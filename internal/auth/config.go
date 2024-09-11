package auth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func MicrosoftConfig() *oauth2.Config {
	MicrosoftLoginConfig := oauth2.Config{
		RedirectURL:  "http://localhost:7331/auth/microsoft_callback",
		ClientID:     os.Getenv("MICROSOFTONLINE_KEY"),
		ClientSecret: os.Getenv("MICROSOFTONLINE_SECRET"),
		Scopes:       []string{"openid", "email", "user.read"},
		Endpoint:     microsoft.AzureADEndpoint("6c7f3fca-54d7-42ee-998d-5a3653dd59bc"),
	}
	MicrosoftLoginConfig.Endpoint.AuthStyle = oauth2.AuthStyleInParams

	return &MicrosoftLoginConfig
}
