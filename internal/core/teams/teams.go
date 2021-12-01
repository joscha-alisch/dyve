package teams

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type Team struct {
	Id           string `json:"id" bson:"id"`
	TeamSettings `bson:",inline"`
}

type TeamSettings struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Access      AccessGroups `json:"access"`
}

type AccessGroups struct {
	Admin  []string `json:"admin"`
	Member []string `json:"member"`
	Viewer []string `json:"viewer"`
}

type TeamPage struct {
	sdk.Pagination
	Teams []Team `json:"teams"`
}

type ByAccess struct {
	Admin  []Team
	Member []Team
	Viewer []Team
}
