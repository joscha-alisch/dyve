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
	j, ok := r.db.AcceptReconcileJob(1 * time.Minute)
	if !ok {
		return false, nil
	}

	log.Info().Interface("job", j).Msg("reconciling")

	var err error
	switch j.Type {
	case ReconcileOrganizations:
		err = r.reconcileOrganizations(j)
	case ReconcileSpaces:
		err = r.reconcileSpaces(j)
	case ReconcileApps:
		err = r.reconcileApps(j)
	}

	return true, err
}

func (r *reconciler) reconcileSpaces(j ReconcileJob) error {
	spaces, err := r.cf.ListSpaces(j.Guid)
	if errors.Is(err, errNotFound) {
		r.db.DeleteOrg(j.Guid)
		return nil
	} else if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	_ = r.db.UpsertOrgSpaces(j.Guid, spaces)
	return nil
}

func (r *reconciler) reconcileApps(j ReconcileJob) error {
	apps, err := r.cf.ListApps(j.Guid)
	if errors.Is(err, errNotFound) {
		r.db.DeleteSpace(j.Guid)
		return nil
	} else if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	_ = r.db.UpsertSpaceApps(j.Guid, apps)
	return nil
}

func (r *reconciler) reconcileOrganizations(j ReconcileJob) error {
	orgs, err := r.cf.ListOrgs()
	if err != nil {
		return err
	}

	_ = r.db.UpsertOrgs(j.Guid, orgs)

	return nil
}
