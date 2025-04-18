package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type GoogleTokenInfo struct {
	Email    string `json:"email"`
	UserID   string `json:"sub"`
	Audience string `json:"aud"`
}

func VerifyGoogleToken(idToken string) (*GoogleTokenInfo, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid ID token")
	}

	var tokenInfo GoogleTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	if err != nil {
		return nil, err
	}

	// OPTIONAL: Validasi audience
	expectedAud := os.Getenv("GOOGLE_CLIENT_ID")
	if tokenInfo.Audience != expectedAud {
		return nil, errors.New("audience mismatch")
	}

	return &tokenInfo, nil
}