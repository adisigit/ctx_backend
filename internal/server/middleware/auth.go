package middleware

import (
	"ctx_backend/internal/auth"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

const ContextKeyUserID = "user_id"

func GetCookie(c huma.Context, name string) (string, error) {
	cookieHeader := c.Header("Cookie")
	if cookieHeader == "" {
		return "", errors.New("no cookie header")
	}
	req := http.Request{Header: http.Header{"Cookie": []string{cookieHeader}}}
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func RequireAuth(jwtService *auth.JWTService, cliTokenService *auth.CLITokenService) func(huma.Context, func(huma.Context)) {
	return func(c huma.Context, next func(huma.Context)) {
		var userID string
		var err error
		if header := c.Header("Authorization"); header != "" {
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				huma.WriteErr(nil, c, http.StatusUnauthorized, "invalid authorization")
				return
			}
			userID, err = cliTokenService.Verify(parts[1])
		} else {
			accessToken, cerr := GetCookie(c, "access_token")
			if cerr != nil || accessToken == "" {
				huma.WriteErr(nil, c, http.StatusUnauthorized, "unauthorized")
				return
			}
			userID, err = jwtService.Verify(accessToken)
		}
		if err != nil || userID == "" {
			huma.WriteErr(nil, c, http.StatusUnauthorized, err.Error())
			return
		}
		newCtx := huma.WithValue(c, ContextKeyUserID, userID)
		next(newCtx)
	}
}
