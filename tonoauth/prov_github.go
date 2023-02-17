package tonoauth

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func NewGithubProvider(opts ProviderOpts) *Provider {
	conf := oauth2.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: "",
		Scopes:      []string{"user:email"},
	}
	getEmail := func(ctx context.Context, token *oauth2.Token) (string, error) {
		client := conf.Client(ctx, token)

		type EmailRes struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		emailsOut := []EmailRes{}

		// Request the user's emails from github, then parse the response
		ezreq := EzHttpRequest{
			Method: "GET",
			Url:    "https://api.github.com/user/emails",
			Headers: gin.H{
				"Accept":               "application/vnd.github+json",
				"X-GitHub-Api-Version": "2022-11-28",
			},
		}
		if err := ezreq.Do(client, &emailsOut); err != nil {
			return "", err
		}

		// Error checking is important
		if len(emailsOut) == 0 {
			return "", ErrNoEmail
		}
		// Try to find the primary email
		for _, em := range emailsOut {
			if em.Primary {
				return em.Email, nil
			}
		}
		// No primary found, fallback to 0th email
		return emailsOut[0].Email, nil
	}
	return &Provider{
		Name:     "github",
		Config:   conf,
		GetEmail: getEmail,
	}
}
