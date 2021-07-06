package cloudfoundry

import "time"

type Database interface {
	AcceptReconcileJob(olderThan time.Duration) (ReconcileJob, bool)

	UpsertOrgs(cfGuid string, orgs []Org) error
	UpsertOrgSpaces(orgGuid string, spaces []Space) error
	UpsertSpaceApps(spaceGuid string, apps []App) error

	DeleteOrg(guid string)
	DeleteSpace(guid string)
	DeleteApp(guid string)

	ListApps() ([]App, error)
}
