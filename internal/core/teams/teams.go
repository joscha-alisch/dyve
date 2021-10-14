package teams

import (
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type TeamPage struct {
	sdk.Pagination
	Teams []Team
}

type Team struct {
	Name    string
	Id      string
	Access  AccessConfig
	Members []sdk.Member
}

type AccessConfig struct {
	AllowGroups []AllowGroupConfig
}

type AllowGroupProvider struct {
	Provider            string
	ProviderDisplayName string
	Groups              []AllowGroupConfig
}

type AllowGroupConfig struct {
	AccessType AccessType
	Group      sdk.Group
}

type ResourceAccessConfig struct {
	Provider            string
	ProviderDisplayName string
	Type                provider.Type
	AllowWith           []AllowResourceConfig
}

type AllowResourceConfig map[string]string

type AccessType string

const (
	AccessTypeAdmin  AccessType = "admin"
	AccessTypeMember AccessType = "member"
	AccessTypeViewer AccessType = "viewer"
)
