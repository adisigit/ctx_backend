package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
)

type OAuthState struct {
	Client    string `json:"client"`
	SessionID string `json:"session_id,omitempty"`
	CSRF      string `json:"csrf"`
}

func EncodeState(s OAuthState) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DecodeState(raw string) (OAuthState, error) {
	var s OAuthState
	b, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(b, &s)
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random string: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}
