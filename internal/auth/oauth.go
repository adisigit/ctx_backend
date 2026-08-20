package auth

import (
	"fmt"

	"golang.org/x/oauth2"
)

type ProviderUser struct {
	ID     string
	Name   string
	Email  string
	Avatar string
}

type Provider interface {
	Config() *oauth2.Config
	GetUser(token *oauth2.Token) (*ProviderUser, error)
}

var providers = map[string]Provider{}

func RegisterProvider(name string, p Provider) {
	providers[name] = p
}

func GetProvider(name string) (Provider, error) {
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return p, nil
}
