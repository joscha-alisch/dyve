package github

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

func NewGroupProvider(db Database) sdk.GroupProvider {
	return &provider{
		db: db,
	}
}

type provider struct {
	db Database
}

func (p *provider) ListGroups() ([]sdk.Group, error) {
	ghTeams, err := p.db.ListTeams()
	if err != nil {
		return nil, err
	}

	var res []sdk.Group
	for _, team := range ghTeams {
		res = append(res, team.toSdkGroup())
	}
	return res, nil
}

func (p *provider) GetGroup(id string) (sdk.Group, error) {
	team, err := p.db.GetTeam(id)
	if err != nil {
		return sdk.Group{}, err
	}

	return team.toSdkGroup(), nil
}

func (t Team) toSdkGroup() sdk.Group {
	var members []sdk.Member
	for _, member := range t.Members {
		members = append(members, sdk.Member{
			Id:   member.Guid,
			Name: member.Name,
		})
	}
	return sdk.Group{
		Id:      t.Guid,
		Name:    t.Name,
		Members: members,
	}
}
