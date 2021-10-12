package github

import (
	"github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

type Database interface {
	AcceptReconcileJob(olderThan time.Duration) (reconciliation.Job, bool)

	ListTeams() ([]Team, error)
	UpsertOrgTeams(org string, teams []Team) error
	GetTeam(guid string) (Team, error)
	UpdateTeamMembers(guid string, members []Member) error
}
