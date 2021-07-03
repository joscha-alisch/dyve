package cloudfoundry

import (
	"errors"
	"time"
)

/**
The reconciler fetches new reconciliation work from the mongoDatabase and updates the corresponding
item via the CloudFoundry API.

It returns true, if there was work to be done and false, if there was no open reconciliation work.
*/
type Reconciler interface {
	Run() (bool, error)
}

func NewReconciler(db Database, cf API) Reconciler {
	return &reconciler{
		db: db,
		cf: cf,
	}
}

type reconciler struct {
	db Database
	cf API
}

func (r *reconciler) Run() (bool, error) {
	j, ok := r.db.AcceptReconcileJob(time.Now(), time.Now())
	if !ok {
		return false, nil
	}

	switch j.Type {
	case ReconcileOrg:
		o, err := r.cf.GetOrg(j.Guid)
		if errors.Is(err, errNotFound) {
			r.db.DeleteOrg(j.Guid)
			return true, nil
		} else if err != nil {
			return true, &errReconcileFailed{Err: err, Job: j}
		}

		_ = r.db.UpsertOrg(o)
	case ReconcileSpace:
		s, apps, err := r.cf.GetSpace(j.Guid)
		if errors.Is(err, errNotFound) {
			r.db.DeleteSpace(j.Guid)
			return true, nil
		} else if err != nil {
			return true, &errReconcileFailed{Err: err, Job: j}
		}

		_ = r.db.UpsertSpace(s)
		_ = r.db.UpsertApps(apps)
	}

	return true, nil
}
