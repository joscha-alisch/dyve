package cloudfoundry

import "time"

type Database interface {
	AcceptReconcileJob(olderThan time.Time, againAt time.Time) (ReconcileJob, bool)
	UpsertOrg(o Org) error
	UpsertSpace(s Space) error
	UpsertApps(apps []App) error

	DeleteOrg(guid string)
	DeleteSpace(guid string)
	DeleteApp(guid string)
}
