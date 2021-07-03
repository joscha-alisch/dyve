package cloudfoundry

import "time"

type Database interface {
	AcceptReconcileJob(olderThan time.Time, againAt time.Time) (ReconcileJob, bool)
	UpsertOrg(o Org) error
	UpsertSpace(s Space) error
	UpsertApp(a App) error
}
