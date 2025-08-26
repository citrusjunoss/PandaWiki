package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

const (
	SessionCacheKey = "_session_store"
	SessionName     = "_pw_auth_session"
)

type AuthGroup struct {
	ID        uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string        `json:"name" gorm:"uniqueIndex;size:100;not null"`
	KbID      string        `gorm:"column:kb_id;not null" json:"kb_id,omitempty"`
	AuthIDs   pq.Int64Array `json:"auth_ids" gorm:"type:int[]"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (AuthGroup) TableName() string {
	return "auth_groups"
}

type AuthInfo struct {
	ID           uint         `gorm:"column:id" json:"id,omitempty"`
	AuthUserInfo AuthUserInfo `json:"auth_user_info" gorm:"type:jsonb"`
}

type AuthUserInfo struct {
	Username  string `json:"username,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Email     string `json:"email,omitempty"`
}

func (s *AuthUserInfo) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("invalid user info type:", value))
	}
	return json.Unmarshal(bytes, s)
}

func (s *AuthUserInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func GetAuthID(c echo.Context) uint {
	userId, ok := c.Get("user_id").(uint)
	if !ok {
		return 0
	}
	return userId
}
