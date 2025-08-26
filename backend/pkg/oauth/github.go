package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/chaitin/panda-wiki/log"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserInfoURL  = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
	githubCallbackPath = "/share/pro/v1/openapi/github/callback"
)

func NewGithubClient(ctx context.Context, logger *log.Logger, clientID, clientSecret, redirectURI string) (*Client, error) {

	redirectURL, _ := url.Parse(redirectURI)
	redirectURL.Path = githubCallbackPath
	redirectURI = redirectURL.String()

	config := Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email"},
		AuthorizeURL: githubAuthorizeURL,
		TokenURL:     githubTokenURL,
		UserInfoURL:  githubUserInfoURL,
		IDField:      "id",
		NameField:    "login",
		AvatarField:  "avatar_url",
		EmailField:   "email",
		RedirectURI:  redirectURI,
	}

	return &Client{
		ctx:    ctx,
		logger: logger.WithModule("pkg.oauth"),
		oauth: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  config.AuthorizeURL,
				TokenURL: config.TokenURL,
			},
			RedirectURL: redirectURI,
			Scopes:      config.Scopes,
		},
		config: &config,
	}, nil
}

func (c *Client) GetGithubPrimaryEmail(token *oauth2.Token) (string, error) {
	client := c.oauth.Client(c.ctx, token)
	type Email struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	resp, err := client.Get(githubUserEmailURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	c.logger.Info("GetGithubPrimaryEmail:", log.Any("buf", string(buf)))

	var emails []Email
	if err := json.Unmarshal(buf, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	return "", errors.New("no primary verified email found")
}
