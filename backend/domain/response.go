package domain

type PWResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Code    int    `json:"code"`
}

type PWResponseErrCode PWResponse

var (
	ErrCodePermissionDenied = PWResponseErrCode{"permission denied", false, nil, 40003}
)
