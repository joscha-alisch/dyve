package cloudfoundry

import (
	"github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

type Database interface {
	AcceptReconcileJob(olderThan time.Duration) (reconciliation.Job, bool)

	UpsertOrgs(cfGuid string, orgs []Org) error
	UpsertOrgSpaces(orgGuid string, spaces []Space) error
	UpsertSpaceApps(spaceGuid string, apps []App) error

	DeleteOrg(guid string)
	DeleteSpace(guid string)
	DeleteApp(guid string)

	ListApps() ([]App, error)
	GetApp(id string) (App, error)

	Cached(id string, duration time.Duration, f func() (interface{}, error)) (interface{}, error)
}
