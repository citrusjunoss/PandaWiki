package dingtalk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcard_1_0 "github.com/alibabacloud-go/dingtalk/card_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/v2/oauth2_1_0"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	callbackPath = "/share/pro/v1/openapi/dingtalk/callback"
	userInfoUrl  = "https://api.dingtalk.com/v1.0/contact/users/me"
)

type Client struct {
	ctx                 context.Context
	clientID            string
	clientSecret        string
	oauthClient         *dingtalkoauth2_1_0.Client
	cardClient          *dingtalkcard_1_0.Client
	dingTalkAuthURL     string
	dingTalkUserInfoURL string
}

// UserInfo 用于解析获取用户信息的接口返回
type UserInfo struct {
	Nick      string `json:"nick"`
	UnionID   string `json:"unionId"`
	OpenID    string `json:"openId"`
	AvatarURL string `json:"avatarUrl"`
	StateCode string `json:"stateCode"`
}

func NewDingTalkClient(ctx context.Context, clientId, clientSecret string) (*Client, error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	oauthClient, err := dingtalkoauth2_1_0.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth client: %w", err)
	}
	cardClient, err := dingtalkcard_1_0.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create card client: %w", err)
	}
	return &Client{
		ctx:                 ctx,
		clientID:            clientId,
		clientSecret:        clientSecret,
		oauthClient:         oauthClient,
		cardClient:          cardClient,
		dingTalkAuthURL:     "https://login.dingtalk.com/oauth2/auth",
		dingTalkUserInfoURL: "https://oapi.dingtalk.com/sns/userinfo",
	}, nil
}

// GenerateAuthURL 生成钉钉授权URL
func (c *Client) GenerateAuthURL(redirectURI string, state string) string {
	redirectURL, _ := url.Parse(redirectURI)
	redirectURL.Path = callbackPath
	redirectURI = redirectURL.String()

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", c.clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "openid")
	params.Add("state", state)
	params.Add("prompt", "consent")

	return fmt.Sprintf("%s?%s", c.dingTalkAuthURL, params.Encode())
}

func (c *Client) GetAccessTokenByCode(code string) (string, error) {
	request := &dingtalkoauth2_1_0.GetUserTokenRequest{
		ClientId:     tea.String(c.clientID),
		ClientSecret: tea.String(c.clientSecret),
		Code:         tea.String(code),
		GrantType:    tea.String("authorization_code"),
	}
	response, err := c.oauthClient.GetUserToken(request)
	if err != nil {
		return "", fmt.Errorf("获取用户access token失败: %w", err)
	}
	accessToken := tea.StringValue(response.Body.AccessToken)
	return accessToken, nil
}

func (c *Client) GetUserInfoByCode(code string) (*UserInfo, error) {

	req, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建GET请求失败: %w", err)
	}
	accessToken, err := c.GetAccessTokenByCode(code)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送GET请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("钉钉API返回错误状态: %s, 响应: %s", resp.Status, string(body))
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析JSON响应失败: %w", err)
	}

	return &userInfo, nil
}
