package apigee

import (
	"context"

	"golang.org/x/oauth2"
)

type TokenSource struct {
	source oauth2.TokenSource
	client oauth2.Config
	opts   *ApigeeClientOptions
}

func NewTokenSource(config *ApigeeClientOptions) *TokenSource {
	c := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL:  config.LoginBaseUrl,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		ClientID:     "edgecli",
		ClientSecret: "edgeclisecret",
	}

	return &TokenSource{source: nil, client: c, opts: config}
}

func (t *TokenSource) newSource(ctx context.Context) error {
	token, err := t.client.PasswordCredentialsToken(ctx, t.opts.Auth.Username, t.opts.Auth.Password)
	if err != nil {
		return err
	}

	t.source = t.client.TokenSource(ctx, token)
	return nil
}

func (t *TokenSource) GetToken(ctx context.Context) (*oauth2.Token, error) {
	if t.source == nil {
		if err := t.newSource(ctx); err != nil {
			return nil, err
		}
	}

	token, err := t.source.Token()
	if err != nil {
		if err := t.newSource(ctx); err != nil {
			return nil, err
		}
		return t.source.Token()
	}

	return token, nil
}
