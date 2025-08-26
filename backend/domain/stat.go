package domain

import "time"

type StatPageScene int

const (
	StatPageSceneWelcome StatPageScene = iota + 1
	StatPageSceneNodeDetail
	StatPageSceneChat
	StatPageSceneLogin
)

type StatPage struct {
	ID          int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	KBID        string        `json:"kb_id"`
	NodeID      string        `json:"node_id"`
	UserID      uint          `json:"user_id"`
	SessionID   string        `json:"session_id"`
	Scene       StatPageScene `json:"scene"` // 1: welcome, 2: detail, 3: chat, 4: login
	IP          string        `json:"ip"`
	UA          string        `json:"ua"`
	BrowserName string        `json:"browser_name"`
	BrowserOS   string        `json:"browser_os"`
	Referer     string        `json:"referer"`
	RefererHost string        `json:"referer_host"`
	CreatedAt   time.Time     `json:"created_at"`
}

type StatPageReq struct {
	Scene  StatPageScene `json:"scene" validate:"required,oneof=1 2 3 4"`
	NodeID string        `json:"node_id"`
}

type HotPageResp struct {
	Scene    StatPageScene `json:"scene"`
	NodeID   string        `json:"node_id"`
	NodeName string        `json:"node_name" gorm:"-"`
	Count    int           `json:"count"`
}

type HotRefererHostResp struct {
	RefererHost string `json:"referer_host"`
	Count       int    `json:"count"`
}

type HotBrowserResp struct {
	OS      []BrowserCount `json:"os"`
	Browser []BrowserCount `json:"browser"`
}

type BrowserCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type StatPageCountResp struct {
	IPCount           int64 `json:"ip_count"`
	SessionCount      int64 `json:"session_count"`
	PageVisitCount    int64 `json:"page_visit_count"`
	ConversationCount int64 `json:"conversation_count"`
}

type InstantCountResp struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

type InstantPageResp struct {
	Scene     StatPageScene `json:"scene"`
	NodeID    string        `json:"node_id"`
	NodeName  string        `json:"node_name" gorm:"-"`
	IP        string        `json:"ip"`
	IPAddress IPAddress     `json:"ip_address" gorm:"-"`
	CreatedAt time.Time     `json:"created_at"`

	UserID uint          `json:"user_id"`
	Info   *AuthUserInfo `json:"info"`
}

type ConversationDistributionResp struct {
	AppType AppType `json:"app_type"`
	AppID   string  `json:"app_id"`
	Count   int     `json:"count"`
}
