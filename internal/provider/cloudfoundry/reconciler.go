package cloudfoundry

import (
	"errors"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

const (
	ReconcileOrganizations recon.Type = "organizations"
	ReconcileSpaces        recon.Type = "spaces"
	ReconcileApps          recon.Type = "apps"
)

/**
The reconciler fetches new reconciliation work from the database and updates the corresponding
item via the CloudFoundry API.

It returns true, if there was work to be done and false, if there was no open reconciliation work.
*/
func NewReconciler(db Database, cf API, olderThan time.Duration) recon.Reconciler {
	if olderThan == 0 {
		olderThan = time.Minute
	}

	r := &reconciler{
		Reconciler: recon.NewReconciler(db, olderThan),
		db:         db,
		cf:         cf,
	}

	r.Handler(ReconcileOrganizations, r.reconcileOrganizations)
	r.Handler(ReconcileSpaces, r.reconcileSpaces)
	r.Handler(ReconcileApps, r.reconcileApps)

	return r
}

type reconciler struct {
	recon.Reconciler

	db Database
	cf API
}

func (r *reconciler) reconcileSpaces(j recon.Job) error {
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

func (r *reconciler) reconcileApps(j recon.Job) error {
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

func (r *reconciler) reconcileOrganizations(j recon.Job) error {
	orgs, err := r.cf.ListOrgs()
	if err != nil {
		return err
	}

	_ = r.db.UpsertOrgs(j.Guid, orgs)

	return nil
}
