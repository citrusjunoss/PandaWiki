package v1

import (
	"time"

	"github.com/chaitin/panda-wiki/domain"
)

type ShareNodeListItemReq struct {
}

type ShareNodeListItemResp struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Type      domain.NodeType `json:"type"`
	ParentID  string          `json:"parent_id"`
	Position  float64         `json:"position"`
	Emoji     string          `json:"emoji"`
	UpdatedAt time.Time       `json:"updated_at"`
}
