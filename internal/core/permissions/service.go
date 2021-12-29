package permissions

import "github.com/joscha-alisch/dyve/internal/core/teams"

type Service interface {
	PermissionsFor(groups []string) (Permissions, error)
}

func NewService(teams teams.Service) Service {
	return &service{
		teams: teams,
	}
}

type service struct {
	teams teams.Service
}

func (s *service) PermissionsFor(groups []string) (Permissions, error) {
	_, err := s.teams.TeamsForGroups(groups)
	if err != nil {
		return Permissions{}, err
	}

	return Permissions{}, nil
}