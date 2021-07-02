package cloudfoundry

/**
The reconciler fetches new reconciliation work from the database and updates the corresponding
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
	j := r.db.FetchReconcileJob()
	if j == nil {
		return false, nil
	}

	switch j.Type {
	case ReconcileOrg:
		o, err := r.cf.GetOrg(j.Guid)
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: *j}
		}

		_ = r.db.UpdateOrg(j.Guid, o)
	case ReconcileSpace:
		s, err := r.cf.GetSpace(j.Guid)
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: *j}
		}

		_ = r.db.UpdateSpace(j.Guid, s)
	case ReconcileApp:
		a, err := r.cf.GetApp(j.Guid)
		if err != nil {
			return true, &errReconcileFailed{Err: err, Job: *j}
		}

		_ = r.db.UpdateApp(j.Guid, a)
	}

	return true, nil
}
