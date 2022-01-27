package fakes

import (
	"github.com/joscha-alisch/dyve/internal/core/teams"
)

type RecordingTeamsService struct {
	Err      error
	Team     teams.Team
	Page     teams.TeamPage
	Teams    []teams.Team
	Record   TeamsRecorder
	ByAccess teams.ByAccess
}

func (a *RecordingTeamsService) EnsureIndices() error {
	//TODO implement me
	panic("implement me")
}

func (a *RecordingTeamsService) ListTeamsPaginated(perPage int, page int) (teams.TeamPage, error) {
	a.Record.PerPage = perPage
	a.Record.Page = page

	if a.Err != nil {
		return teams.TeamPage{}, a.Err
	}
	return a.Page, nil
}

func (a *RecordingTeamsService) GetTeam(id string) (teams.Team, error) {
	a.Record.TeamId = id
	if a.Err != nil {
		return teams.Team{}, a.Err
	}
	return a.Team, nil
}

func (a *RecordingTeamsService) DeleteTeam(id string) error {
	a.Record.TeamId = id
	return a.Err
}

func (a *RecordingTeamsService) CreateTeam(id string, data teams.TeamSettings) error {
	a.Record.TeamId = id
	a.Record.TeamData = data
	return a.Err
}

func (a *RecordingTeamsService) UpdateTeam(id string, data teams.TeamSettings) error {
	a.Record.TeamId = id
	a.Record.TeamData = data
	return a.Err
}

func (a *RecordingTeamsService) TeamsForGroups(groups []string) (teams.ByAccess, error) {
	a.Record.Groups = groups
	if a.Err != nil {
		return teams.ByAccess{}, a.Err
	}
	return a.ByAccess, nil
}

type TeamsRecorder struct {
	Team     teams.Team
	PerPage  int
	Page     int
	TeamId   string
	Teams    []teams.Team
	TeamData teams.TeamSettings
	Groups   []string
}
