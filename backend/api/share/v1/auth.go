package v1

import "github.com/chaitin/panda-wiki/consts"

type AuthLoginSimpleReq struct {
	Password string `json:"password" validate:"required"`
}

type AuthLoginSimpleResp struct {
}

type AuthGetReq struct {
}
type AuthGetResp struct {
	AuthType   consts.AuthType   `json:"auth_type"`
	SourceType consts.SourceType `json:"source_type"`
}
