package cloudfoundry

import "time"

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
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: j}
		}

		_ = r.db.UpsertOrg(o)
		for _, space := range o.Spaces {
			_ = r.db.UpsertJob(ReconcileJob{
				Type: ReconcileSpace,
				Guid: space,
			})
		}

	case ReconcileSpace:
		s, err := r.cf.GetSpace(j.Guid)
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: j}
		}

		_ = r.db.UpsertSpace(s)
	case ReconcileApp:
		a, err := r.cf.GetApp(j.Guid)
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: j}
		}

		_ = r.db.UpsertApp(a)
	}

	return true, nil
}
