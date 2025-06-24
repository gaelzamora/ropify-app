package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	FacebookOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	// Para Twitter, necesitarás usar OAuth1 o su API OAuth2 más reciente
	// La implementación difiere un poco, aquí un placeholder:
	TwitterOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("TWITTER_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("TWITTER_REDIRECT_URL"),
		Scopes:       []string{"tweet.read", "users.read", "offline.access"},
		// Tendrías que implementar el endpoint para Twitter
	}
)
