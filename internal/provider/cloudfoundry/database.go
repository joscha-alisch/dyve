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

	DeleteOrg(guid string) (bool, error)
	DeleteSpace(guid string) (bool, error)
	DeleteApp(guid string) (bool, error)

	ListApps() ([]App, error)
	GetApp(id string) (App, error)

	Cached(id string, duration time.Duration, cached interface{}, f func() (interface{}, error)) (interface{}, error)
}
