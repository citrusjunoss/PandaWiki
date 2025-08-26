package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/chaitin/panda-wiki/log"
)

const (
	// AuthURL api doc https://developer.work.weixin.qq.com/document/path/98152
	AuthURL     = "https://login.work.weixin.qq.com/wwlogin/sso/login"
	TokenURL    = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	UserInfoURL = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"

	UserDetailURL = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	callbackPath  = "/share/pro/v1/openapi/wecom/callback"
)

var oauthEndpoint = oauth2.Endpoint{
	AuthURL:  AuthURL,
	TokenURL: TokenURL,
}

// Client 企业微信客户端
type Client struct {
	context     context.Context
	oauthConfig *oauth2.Config
	logger      *log.Logger
	corpID      string
	agentID     string
}

type TokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserInfoResponse struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	UserID         string `json:"userid"`
	UserTicket     string `json:"user_ticket"`
	OpenID         string `json:"openid"`
	ExternalUserid string `json:"external_userid"`
}

type UserDetailResponse struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Userid     string `json:"userid"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	OpenUserid string `json:"open_userid"`
}

func NewClient(ctx context.Context, logger *log.Logger, corpID, corpSecret, agentID, redirectURI string) (*Client, error) {
	redirectURL, _ := url.Parse(redirectURI)
	redirectURL.Path = callbackPath
	redirectURI = redirectURL.String()

	oauthConfig := &oauth2.Config{
		ClientID:     corpID,
		ClientSecret: corpSecret,
		RedirectURL:  redirectURI,
		Endpoint:     oauthEndpoint,
		Scopes:       []string{"snsapi_privateinfo"},
	}

	return &Client{
		context:     ctx,
		logger:      logger.WithModule("wecom.client"),
		oauthConfig: oauthConfig,
		corpID:      corpID,
		agentID:     agentID,
	}, nil
}

// GenerateAuthURL 生成授权 URL
func (c *Client) GenerateAuthURL(state string) string {
	params := url.Values{}
	params.Set("appid", c.corpID)
	params.Set("redirect_uri", c.oauthConfig.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_privateinfo")
	params.Set("login_type", "CorpApp")
	params.Set("agentid", c.agentID)
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", AuthURL, params.Encode())
}

// GetAccessToken 获取企业微信访问令牌
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	params := url.Values{}
	params.Set("corpid", c.corpID)
	params.Set("corpsecret", c.oauthConfig.ClientSecret)

	resp, err := http.Get(fmt.Sprintf("%s?%s", TokenURL, params.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("failed to get access token: %s", tokenResp.ErrMsg)
	}

	return tokenResp.AccessToken, nil
}

// GetUserInfoByCode 通过授权码获取用户信息
func (c *Client) GetUserInfoByCode(ctx context.Context, code string) (*UserDetailResponse, error) {

	accessToken, err := c.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("code", code)

	userInfoURL := fmt.Sprintf("%s?%s", UserInfoURL, params.Encode())

	c.logger.Info("GetUserInfoByCode", log.Any("userInfoURL", userInfoURL))

	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	rawBody, _ := io.ReadAll(resp.Body)
	c.logger.Info("GetUserInfoByCode raw resp:", log.Any("raw", string(rawBody)))

	resp.Body = io.NopCloser(bytes.NewReader(rawBody))

	var userInfoResp UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResp); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}

	c.logger.Info("GetUserInfoByCode resp:", log.Any("resp", userInfoResp))

	if userInfoResp.ErrCode != 0 {
		return nil, fmt.Errorf("failed to get user info: %s", userInfoResp.ErrMsg)
	}

	detailParams := url.Values{}
	detailParams.Set("access_token", accessToken)
	detailParams.Set("userid", userInfoResp.UserID)

	userDetailURL := fmt.Sprintf("%s?%s", UserDetailURL, detailParams.Encode())

	detailResp, err := http.Get(userDetailURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user detail: %w", err)
	}
	defer detailResp.Body.Close()

	var UserDetailResp UserDetailResponse
	if err := json.NewDecoder(detailResp.Body).Decode(&UserDetailResp); err != nil {
		return nil, fmt.Errorf("failed to decode user detail response: %w", err)
	}

	c.logger.Info("GetUserInfoByCode detail info", log.Any("resp", UserDetailResp))

	if UserDetailResp.Errcode != 0 {
		return nil, fmt.Errorf("failed to get user detail: %s", UserDetailResp.Errmsg)
	}

	return &UserDetailResp, nil
}
