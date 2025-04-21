package utils

import (
	"context"
	"errors"
	"os"

	"google.golang.org/api/idtoken"
)

type GoogleUserInfo struct {
	Email     string
	Name      string
	Picture   string
	GoogleID  string
	Audience  string
}

func VerifyGoogleToken(idToken string) (*GoogleUserInfo, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token")
	}

	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	googleID, _ := payload.Claims["sub"].(string)
	aud, _ := payload.Claims["aud"].(string)

	return &GoogleUserInfo{
		Email:    email,
		Name:     name,
		Picture:  picture,
		GoogleID: googleID,
		Audience: aud,
	}, nil
}