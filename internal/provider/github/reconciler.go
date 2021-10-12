package github

import (
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

const (
	ReconcileTeams   recon.Type = "teams"
	ReconcileMembers recon.Type = "members"
)

/**
The reconciler fetches new reconciliation work from the database and updates the corresponding
item via the github API.

It returns true, if there was work to be done and false, if there was no open reconciliation work.
*/
func NewReconciler(db Database, gh API, olderThan time.Duration) recon.Reconciler {
	if olderThan == 0 {
		olderThan = time.Minute
	}

	r := &reconciler{
		Reconciler: recon.NewReconciler(db, olderThan),
		db:         db,
		gh:         gh,
	}

	r.Handler(ReconcileTeams, r.reconcileTeams)
	r.Handler(ReconcileMembers, r.reconcileMembers)

	return r
}

type reconciler struct {
	recon.Reconciler

	gh API
	db Database
}

func (r *reconciler) reconcileTeams(j recon.Job) error {
	teams, err := r.gh.ListTeams(j.Guid)
	if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}
	_ = r.db.UpsertOrgTeams(j.Guid, teams)
	return nil
}

func (r *reconciler) reconcileMembers(j recon.Job) error {
	t, err := r.db.GetTeam(j.Guid)
	if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	members, err := r.gh.ListMembers(t.Org.Guid, t.Slug)
	if err != nil {
		return &errReconcileFailed{Err: err, Job: j}
	}

	err = r.db.UpdateTeamMembers(t.Guid, members)
	if err != nil {
		return err
	}

	return nil
}
