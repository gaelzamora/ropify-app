package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gaelzamora/ropify-app/config"
	"github.com/gaelzamora/ropify-app/models"
	"github.com/gaelzamora/ropify-app/utils"
	"github.com/golang-jwt/jwt/v5"
)

type OAuthService struct {
	repository models.AuthRepository
}

func (s *OAuthService) HandleGoogleLogin(ctx context.Context, code string) (string, *models.User, error) {
	// Exchange code for token
	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return "", nil, fmt.Errorf("code exchange failed: %v", err)
	}

	// Get user info from Google
	userInfo, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return "", nil, err
	}

	// Check if user exists or create new user
	user, err := s.repository.GetUser(ctx, "email = ?", userInfo.Email)
	if err != nil {
		// Create new user
		newUser := &models.User{
			Email:     userInfo.Email,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Username:  userInfo.Email, // O generar un nombre de usuario único
			GoogleID:  &userInfo.ID,
		}

		user, err = s.repository.RegisterUser(ctx, &models.AuthCredentials{
			Email:     newUser.Email,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Username:  newUser.Username,
			// Generar contraseña aleatoria o dejar en blanco
		})

		if err != nil {
			return "", nil, err
		}
	} else {
		// Actualizar GoogleID si es necesario
		if user.GoogleID == nil || *user.GoogleID != userInfo.ID {
			// Actualizar el usuario con el GoogleID
			// Implementar esta lógica en tu repositorio
		}
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 168).Unix(),
	}

	jwtToken, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", nil, err
	}

	return jwtToken, user, nil
}

// Estructura para la respuesta de Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func (s *OAuthService) getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// Implementar métodos similares para Facebook y Twitter

func NewOAuthService(repository models.AuthRepository) *OAuthService {
	return &OAuthService{
		repository: repository,
	}
}
