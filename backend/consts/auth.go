package consts

type SourceType string

const (
	SourceTypeDingTalk SourceType = "dingtalk"
	SourceTypeFeishu   SourceType = "feishu"
	SourceTypeWeCom    SourceType = "wecom"
	SourceTypeOAuth    SourceType = "oauth"
	SourceTypeGitHub   SourceType = "github"
	SourceTypeCAS      SourceType = "cas"
	SourceTypeLDAP     SourceType = "ldap"
)

type AuthType string

const (
	AuthTypeNull       AuthType = ""           // 无认证
	AuthTypeSimple     AuthType = "simple"     // 简单口令
	AuthTypeEnterprise AuthType = "enterprise" // 企业认证

)
