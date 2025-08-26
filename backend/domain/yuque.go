package domain

type YuqueReq struct {
	KbID string `json:"kb_id" binding:"required" validate:"required"`
}

type YuqueResp struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}
