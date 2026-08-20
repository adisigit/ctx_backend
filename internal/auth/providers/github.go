package providers

import (
	"context"
	"ctx_backend/internal/auth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubProvider struct{}

func (p *GithubProvider) Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: github.Endpoint,
	}
}

func (p *GithubProvider) GetUser(token *oauth2.Token) (*auth.ProviderUser, error) {
	client := p.Config().Client(context.Background(), token)
	raw, err := fetchGithubUser(client)
	if err != nil {
		return nil, err
	}
	email := raw.Email
	if email == "" {
		email, err = fetchGithubPrimaryEmail(client)
		if err != nil {
			return nil, err
		}
	}
	return &auth.ProviderUser{
		ID:     fmt.Sprintf("%d", raw.ID),
		Name:   raw.Name,
		Email:  email,
		Avatar: raw.AvatarURL,
	}, nil
}

type githubUserResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func fetchGithubUser(client *http.Client) (*githubUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ctx_backend")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var raw githubUserResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func fetchGithubPrimaryEmail(client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ctx_backend")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}
	return "", nil
}
