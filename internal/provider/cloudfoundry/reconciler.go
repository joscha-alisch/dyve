package cloudfoundry

import (
	"errors"
	"github.com/rs/zerolog/log"
	"time"
)

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
	j, ok := r.db.AcceptReconcileJob(20 * time.Second)
	if !ok {
		return false, nil
	}

	log.Info().Interface("job", j).Msg("reconciling")

	var err error
	switch j.Type {
	case ReconcileCF:
		err = r.reconcileCF(j)
	case ReconcileOrg:
		err = r.reconcileOrg(j)
	case ReconcileSpace:
		err = r.reconcileSpace(j)
	}

	return true, err
}

func (r *reconciler) reconcileOrg(j ReconcileJob) error {
	o, err := r.cf.GetOrg(j.Guid)
	if errors.Is(err, errNotFound) {
		r.db.DeleteOrg(j.Guid)
		return nil
	} else if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	_ = r.db.UpsertOrg(o)
	return nil
}

func (r *reconciler) reconcileSpace(j ReconcileJob) error {
	s, apps, err := r.cf.GetSpace(j.Guid)
	if errors.Is(err, errNotFound) {
		r.db.DeleteSpace(j.Guid)
		return nil
	} else if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	_ = r.db.UpsertSpace(s)
	_ = r.db.UpsertApps(apps)
	return nil
}

func (r *reconciler) reconcileCF(j ReconcileJob) error {
	i, err := r.cf.GetCFInfo()
	if err != nil {
		return err
	}

	_ = r.db.UpsertCfInfo(i)

	return nil
}
